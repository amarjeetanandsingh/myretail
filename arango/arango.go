package arango

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type db struct {
	client     driver.Client
	dbInstance driver.Database
}

type DbServer interface {
	FindByID(coll string, key string, result interface{}) error
	Query(string, map[string]interface{}, interface{}) error
	UpdateDoc(coll string, key string, fields interface{}) error
}

func (store *db) FindByID(collection string, key string, result interface{}) error {
	q := `RETURN DOCUMENT(CONCAT(@@coll/@key))`
	bindVar := map[string]interface{}{
		"@coll": collection,
		"key":   key,
	}
	return store.Query(q, bindVar, result)
}

func (store *db) Query(q string, bindVars map[string]interface{}, result interface{}) error {
	ctx := context.Background()
	cur, err := store.dbInstance.Query(ctx, q, bindVars)
	if err != nil {
		return err
	}
	defer cur.Close()

	if _, err := cur.ReadDocument(ctx, result); err != nil {
		return fmt.Errorf("error unmarshling result: %w", err)
	}

	return nil
}

func (store *db) UpdateDoc(coll string, key string, fields interface{}) error {
	ctx := context.Background()
	collection, err := store.dbInstance.Collection(ctx, coll)
	if err != nil {
		return err
	}

	if _, err := collection.UpdateDocument(ctx, key, fields); err != nil {
		return err
	}

	return nil
}

func New(host, port, userName, password string) (DbServer, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://" + host + ":" + port},
	})
	if err != nil {
		return nil, fmt.Errorf("Error connecting ArangoDB %w", err)
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(userName, password),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating arango client: %w", err)
	}

	dbInstance, err := client.Database(nil, "rb-wheels") //TODO db name is product
	if err != nil {
		return nil, fmt.Errorf("error creating db instance: %w", err)
	}

	return &db{
		dbInstance: dbInstance,
		client:     client,
	}, nil
}
