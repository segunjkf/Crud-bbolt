package bolt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"

	"github.com/segunjkf/server/pkg/database"
)

// Bolt is the Bolt database.
// It satisfies the Database Interface.
type Bolt struct{
	db *bolt.DB
}

const ( 
	dbName = "text.db"
	bucketName = "users"
)

// New Return a new Bolt Implementation
func New(ctx context.Context, directory string) (*Bolt, error) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(fmt.Sprintf("%s/%s", directory, dbName), 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Ensure that the bucket exists.
	err  = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(([]byte(bucketName)))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Bolt{
		db: db,
	}, nil
}

// Close closes the database
// Make sure to close the database
func (b *Bolt) Close(ctx context.Context) {
	b.db.Close()
}

type userInfo struct {
	Email string `json:"email"`
	Age   int	 `json:"age"`
}

func (b *Bolt) Create(ctx context.Context, user database.User) error {

	v, err := json.Marshal(userInfo{
		Email: user.Email,
		Age: user.Age,
	})
	if err != nil {
		return err
	}

	b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(user.Name), v)
		return err
	})
	
	return nil
}

//  Get implemetation of the Database interface
func(b *Bolt) GetUser(ctx context.Context, name string) (user *database.User) {

	var raw []byte
	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		raw = b.Get([]byte(name))
		fmt.Printf("The answer is: %s\n", raw)
		return nil
	})
	if len(raw) == 0 {
		return 
	}
	var u database.User

	err := json.Unmarshal(raw, &u)
	if err != nil {
		log.Fatalf("Database corruption: %v", err)
	}
	user = &u
	return 
}

// update implements the database interface
func(b *Bolt) Update(ctx context.Context, name string) (database.User, error) {
	return database.User{}, nil
}