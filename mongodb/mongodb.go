package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Storage interface that is implemented by storage providers
type Storage struct {
	db *mongo.Collection
}

type MongoStorage struct {
	ObjectID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Key      string             `json:"key,omitempty" bson:"key"`
	Value    []byte             `json:"value" bson:"value"`
	Exp      time.Time          `json:"exp" bson:"exp,omitempty"`
}

// New creates a new MongoDB storage
func New(config ...Config) *Storage {
	// Set default config
	cfg := ConfigDefault

	// Override config if provided
	if len(config) > 0 {
		cfg = configDefault(config[0])
	}

	// Set mongo options
	opt := options.Client()
	opt.SetAppName(cfg.AppName)
	opt.SetAuth(cfg.Auth)
	opt.SetAutoEncryptionOptions(cfg.AutoEncryptionOptions)
	opt.SetConnectTimeout(cfg.ConnectTimeout)
	opt.SetCompressors(cfg.Compressors)
	opt.SetDialer(cfg.Dialer)
	opt.SetDirect(cfg.Direct)
	opt.SetDisableOCSPEndpointCheck(cfg.DisableOCSPEndpointCheck)
	opt.SetHeartbeatInterval(cfg.HeartbeatInterval)
	opt.SetHosts(cfg.Hosts)
	opt.SetLocalThreshold(cfg.LocalThreshold)
	opt.SetMaxConnIdleTime(cfg.MaxConnIdleTime)
	opt.SetMaxPoolSize(cfg.MaxPoolSize)
	opt.SetMinPoolSize(cfg.MinPoolSize)
	opt.SetPoolMonitor(cfg.PoolMonitor)
	opt.SetMonitor(cfg.Monitor)
	opt.SetReadConcern(cfg.ReadConcern)
	opt.SetReadPreference(cfg.ReadPreference)
	opt.SetRegistry(cfg.Registry)
	opt.SetReplicaSet(cfg.ReplicaSet)
	opt.SetRetryReads(cfg.RetryReads)
	opt.SetRetryWrites(cfg.RetryWrites)
	opt.SetServerSelectionTimeout(cfg.ServerSelectionTimeout)
	opt.SetSocketTimeout(cfg.SocketTimeout)
	opt.SetTLSConfig(cfg.TLSConfig)
	opt.SetWriteConcern(cfg.WriteConcern)
	opt.SetZlibLevel(cfg.ZlibLevel)
	opt.SetZstdLevel(cfg.ZstdLevel)

	// Create mongo client
	client, err := mongo.NewClient(opt.ApplyURI("mongodb://" + cfg.Addr))
	if err != nil {
		panic(err)
	}

	// Get collection from database
	db := client.Database(cfg.Database).Collection(cfg.Collection)

	// expired data may exist for some time beyond the 60 second period between runs of the background task.
	// more on https://docs.mongodb.com/manual/core/index-ttl/
	indexModel := mongo.IndexModel{
		Keys: bson.D{{"exp", 1}},
		// setting to 0
		// means that documents will remain in the collection
		// until they're explicitly deleted or the collection is dropped.
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	if _, err := db.Indexes().CreateOne(context.TODO(), indexModel); err != nil {
		panic(err)
	}

	return &Storage{
		db: db,
	}
}

// Get value by key
func (s *Storage) Get(key string) ([]byte, error) {
	res := s.db.FindOne(context.TODO(), bson.M{"key": key})
	result := MongoStorage{}

	if err := res.Err(); err != nil {
		return []byte{}, err
	}
	if err := res.Decode(&result); err != nil {
		return []byte{}, err
	}

	return result.Value, nil
}

// Set key with value, replace if document exits
func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
	filter := bson.M{"id": key}
	replace := MongoStorage{
		Key:   key,
		Value: val,
	}

	if exp != 0 {
		replace.Exp = time.Now().Add(exp).UTC()
	}
	_, err := s.db.ReplaceOne(context.TODO(), filter, replace, options.Replace().SetUpsert(true))
	return err
}

// Delete document by key
func (s *Storage) Delete(key string) error {
	_, err := s.db.DeleteOne(context.TODO(), bson.M{"key": key})
	return err
}

// Clear all keys by drop collection
func (s *Storage) Clear() error {
	return s.db.Drop(context.TODO())
}