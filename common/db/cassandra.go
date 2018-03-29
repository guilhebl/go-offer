package db

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/guilhebl/go-offer/common/model"
	"github.com/guilhebl/go-offer/common/util"
	"log"
	"sync"
	"time"
)

// represents a Cassandra driver client
type CassandraClient struct {
	ClusterConfig *gocql.ClusterConfig
}

var instance *CassandraClient
var once sync.Once

// builds a new cassandra db client instance using default port 9042
func BuildInstance(host, username, password, keyspace string, port int) *CassandraClient {
	once.Do(func() {
		cluster := gocql.NewCluster(host)
		cluster.Keyspace = keyspace
		cluster.Port = port
		cluster.ProtoVersion = 4
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: username,
			Password: password,
		}

		instance = &CassandraClient{
			ClusterConfig: cluster,
		}
	})
	return instance
}

func GetInstance() *CassandraClient {
	return instance
}

// gets all offers
func GetOffers() ([]model.Offer, error) {
	log.Printf("Cassandra.GetOffers")

	session, err := GetInstance().ClusterConfig.CreateSession()
	defer session.Close()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var id, externalId, upc, name, partyName, semanticName, mainImageFileUrl, partyImageFileUrl, productCategory string
	var price, rating float32
	var numReviews int
	var created time.Time

	// output list
	list := make([]model.Offer, 0)

	// list all
	selectStatement := `SELECT id, external_id, upc, name, party_name, semantic_name, main_image_file_url, party_image_file_url, product_category, price, rating, num_reviews, created FROM offer`
	iter := session.Query(selectStatement).Iter()
	for iter.Scan(&id, &externalId, &upc, &name, &partyName, &semanticName, &mainImageFileUrl, &partyImageFileUrl, &productCategory, &price, &rating, &numReviews, &created) {
		o := model.NewOffer(id, externalId, upc, name, partyName, semanticName, mainImageFileUrl, partyImageFileUrl, productCategory, price, rating, numReviews, created)
		list = append(list, *o)
	}
	return list, nil
}

// Insert Offer
func InsertOffer(o *model.Offer) (*model.Offer, error) {
	log.Print("Cassandra.InsertOffer")

	session, err := GetInstance().ClusterConfig.CreateSession()
	defer session.Close()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// create new UUID
	o.Id = util.GenerateStringUUID()

	if err := insertOffer(session, o); err != nil {
		return nil, err
	}
	return o, nil
}

// insert Offer
func insertOffer(session *gocql.Session, o *model.Offer) error {
	insertStatement := `
INSERT INTO offer (id, external_id, upc, name, party_name, semantic_name, main_image_file_url, party_image_file_url, product_category, price, rating, num_reviews, created)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// insert an offer
	if err := session.Query(insertStatement,
		o.Id, o.ExternalId, o.Upc, o.Name, o.PartyName, o.SemanticName, o.MainImageFileUrl, o.PartyImageFileUrl, o.ProductCategory, o.Price, o.Rating, o.NumReviews, o.Created).Exec(); err != nil {
		return err
	}
	return nil
}

// Resets DB
func Reset() error {
	log.Print("Cassandra.Reset")

	session, err := GetInstance().ClusterConfig.CreateSession()
	defer session.Close()
	if err != nil {
		log.Print(err)
		return err
	}

	// drop table if exists
	keyspace := GetInstance().ClusterConfig.Keyspace
	dropTable := fmt.Sprintf("DROP TABLE IF EXISTS %s.offer", keyspace)
	if err := session.Query(dropTable).Exec(); err != nil {
		log.Print(err)
		return err
	}

	// create table
	createTableStatement := fmt.Sprintf(`
CREATE TABLE %s.offer (
	id uuid PRIMARY KEY,
	external_id text, 
	upc text, 
	name text, 
	party_name text, 
	semantic_name text, 
	main_image_file_url text, 
	party_image_file_url text, 
	product_category text, 
	price float, 
	rating float, 
	num_reviews int,
	created date);
`, keyspace)

	if err := session.Query(createTableStatement).Exec(); err != nil {
		log.Print(err)
		return err
	}

	//// create index on UPC
	//createIndexUpcStatement := fmt.Sprintf(`
	//CREATE INDEX IF NOT EXISTS offer_upc
	//ON %s.offer (upc)`, keyspace)
	//
	//if err := session.Query(createIndexUpcStatement).Exec(); err != nil {
	//	log.Print(err)
	//	return err
	//}

	// insert sample offers
	if err := insertOffer(session, model.NewOffer(
		util.GenerateStringUUID(), "1", "upc12345678", "offer 1", "amazon.com", "https://amazon.com/offer/001", "https://amazon.com/img/offer/001", "amazon-logo.jpg", "offers", 50.00, 2.5, 50, time.Now(),
	)); err != nil {
		log.Print(err)
		return err
	}

	if err := insertOffer(session, model.NewOffer(
		util.GenerateStringUUID(), "2", "upc22345678", "offer 2", "bestbuy.com", "https://bestbuy.com/offer/001", "https://bestbuy.com/img/offer/001", "bestbuy-logo.jpg", "offers", 60.00, 2.8, 30, time.Now(),
	)); err != nil {
		log.Print(err)
		return err
	}

	if err := insertOffer(session, model.NewOffer(
		util.GenerateStringUUID(), "3", "upc32345678", "offer 3", "walmart.com", "https://walmart.com/offer/001", "https://walmart.com/img/offer/001", "walmart-logo.jpg", "offers", 65.00, 4.5, 60, time.Now(),
	)); err != nil {
		log.Print(err)
		return err
	}

	if err := insertOffer(session, model.NewOffer(
		util.GenerateStringUUID(), "4", "upc42345678", "offer 4", "ebay.com", "https://ebay.com/offer/001", "https://ebay.com/img/offer/001", "ebay-logo.jpg", "offers", 105.00, 3.5, 60, time.Now(),
	)); err != nil {
		log.Print(err)
		return err
	}
	return nil
}
