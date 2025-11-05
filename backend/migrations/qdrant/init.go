package main

import (
	"context"
	"log"
	"os"

	"github.com/qdrant/go-client/qdrant"
)

func main() {
	// Get Qdrant URL from environment
	qdrantURL := os.Getenv("QDRANT_URL")
	if qdrantURL == "" {
		qdrantURL = "localhost:6333"
	}

	// Create Qdrant client
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: qdrantURL,
	})
	if err != nil {
		log.Fatalf("Failed to create Qdrant client: %v", err)
	}

	ctx := context.Background()

	log.Println("🔮 Initializing Qdrant collections...")

	// Create profile embeddings collection
	err = createProfileEmbeddingsCollection(ctx, client)
	if err != nil {
		log.Fatalf("Failed to create profile embeddings collection: %v", err)
	}

	log.Println("✅ Qdrant collections initialized successfully!")
}

func createProfileEmbeddingsCollection(ctx context.Context, client *qdrant.Client) error {
	collectionName := "profile_embeddings"
	
	log.Printf("Creating collection: %s", collectionName)

	// Create collection with vector configuration
	err := client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: &qdrant.VectorsConfig{
			Params: &qdrant.VectorParams{
				Size:     512, // 512-dimensional embeddings
				Distance: qdrant.Distance_Cosine,
			},
		},
		OptimizersConfig: &qdrant.OptimizersConfigDiff{
			IndexingThreshold: qdrant.PtrOf(uint64(10000)),
		},
		HnswConfig: &qdrant.HnswConfigDiff{
			M:                 qdrant.PtrOf(uint64(16)),
			EfConstruct:       qdrant.PtrOf(uint64(100)),
			FullScanThreshold: qdrant.PtrOf(uint64(10000)),
		},
	})

	if err != nil {
		// Collection might already exist, log and continue
		log.Printf("Note: %v (this is OK if collection already exists)", err)
	}

	// Create payload indexes for filtering
	indexes := []struct {
		field string
		ftype qdrant.FieldType
	}{
		{"user_id", qdrant.FieldType_FieldTypeKeyword},
		{"gender", qdrant.FieldType_FieldTypeKeyword},
		{"age", qdrant.FieldType_FieldTypeInteger},
		{"age_group", qdrant.FieldType_FieldTypeKeyword},
		{"active", qdrant.FieldType_FieldTypeBool},
		{"verified", qdrant.FieldType_FieldTypeBool},
		{"location", qdrant.FieldType_FieldTypeGeo},
	}

	for _, idx := range indexes {
		err = client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
			CollectionName: collectionName,
			FieldName:      idx.field,
			FieldType:      idx.ftype,
		})
		if err != nil {
			log.Printf("Warning: Failed to create index for %s: %v", idx.field, err)
		} else {
			log.Printf("  ✓ Created index: %s", idx.field)
		}
	}

	return nil
}
