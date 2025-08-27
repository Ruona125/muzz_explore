package test

import (
	"context"
	"os"
	"testing"

	"explore_service/internal/server"
	"explore_service/internal/storage"
	explorepb "explore_service/proto"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

// startPostgres spins up a temporary PostgreSQL container for testing.
func startPostgres(ctx context.Context, t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	// If TEST_PG_DSN is set, use that external database instead of starting
	// a testcontainer. This makes it easy to point tests at a local
	// docker-compose Postgres for debugging.
	if dsn := os.Getenv("TEST_PG_DSN"); dsn != "" {
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			t.Fatalf("failed to create pgx pool from TEST_PG_DSN: %v", err)
		}
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			t.Fatalf("failed to ping external Postgres from TEST_PG_DSN: %v", err)
		}
		cleanup := func() { pool.Close() }
		return pool, cleanup
	}

	// Create a new Postgres container with default credentials.
	container, err := tcpostgres.RunContainer(ctx, testcontainers.WithImage("postgres:15"))
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	// Build the DSN from the container's connection string.
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}
	// Connect using pgxpool.
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}
	// Return a cleanup function that will close the pool and terminate
	// the container.
	cleanup := func() {
		pool.Close()
		_ = container.Terminate(ctx)
	}
	return pool, cleanup
}

// TestExploreServer exercises the core functionality of the ExploreService.
func TestExploreServer(t *testing.T) {
	ctx := context.Background()
	pool, cleanup := startPostgres(ctx, t)
	defer cleanup()
	// Initialise store and server.
	store, err := storage.NewStore(ctx, pool)
	if err != nil {
		t.Fatalf("failed to initialise store: %v", err)
	}
	srv := server.NewExploreServer(store, 10)
	// Helper to call PutDecision.
	put := func(actor, recipient string, like bool) bool {
		resp, err := srv.PutDecision(ctx, &explorepb.PutDecisionRequest{
			ActorUserId:     actor,
			RecipientUserId: recipient,
			LikedRecipient:  like,
		})
		if err != nil {
			t.Fatalf("PutDecision returned error: %v", err)
		}
		return resp.GetMutualLikes()
	}
	// actor1 likes recipient1 – no mutual like yet.
	if got := put("actor1", "user1", true); got {
		t.Errorf("expected no mutual like on first like")
	}
	// recipient likes actor – now mutual.
	if got := put("user1", "actor1", true); !got {
		t.Errorf("expected mutual like after reciprocal like")
	}
	// actor1 passes user1 – mutual should be false.
	if got := put("actor1", "user1", false); got {
		t.Errorf("expected no mutual like after pass")
	}
	// actor2 and actor3 like user1; actor1 re-likes user1.
	put("actor2", "user1", true)
	put("actor3", "user1", true)
	put("actor1", "user1", true)
	// Count should be 3 (actor1, actor2, actor3 all like user1).
	countResp, err := srv.CountLikedYou(ctx, &explorepb.CountLikedYouRequest{RecipientUserId: "user1"})
	if err != nil {
		t.Fatalf("CountLikedYou returned error: %v", err)
	}
	if countResp.GetCount() != 3 {
		t.Errorf("expected count 3, got %d", countResp.GetCount())
	}
	// List all likers – expect three entries.  Order is by updated_at
	// descending; actor1 is last updated.
	listResp, err := srv.ListLikedYou(ctx, &explorepb.ListLikedYouRequest{RecipientUserId: "user1"})
	if err != nil {
		t.Fatalf("ListLikedYou returned error: %v", err)
	}
	if len(listResp.GetLikers()) != 3 {
		t.Fatalf("expected 3 likers, got %d", len(listResp.GetLikers()))
	}
	// List new likers – actor1 should not be included because user1
	// has liked them back (mutual like).  Expect actors 2 and 3.
	newResp, err := srv.ListNewLikedYou(ctx, &explorepb.ListLikedYouRequest{RecipientUserId: "user1"})
	if err != nil {
		t.Fatalf("ListNewLikedYou returned error: %v", err)
	}
	if len(newResp.GetLikers()) != 2 {
		t.Fatalf("expected 2 new likers, got %d", len(newResp.GetLikers()))
	}
	// Check that actor2 and actor3 are present.
	found := make(map[string]bool)
	for _, l := range newResp.GetLikers() {
		found[l.GetActorId()] = true
	}
	for _, id := range []string{"actor2", "actor3"} {
		if !found[id] {
			t.Errorf("expected to find %s in new likers", id)
		}
	}
}
