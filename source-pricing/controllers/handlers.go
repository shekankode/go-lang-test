package controllers



import (
	"fmt"
	"net/http"
	"github.com/go-redis/redis"
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	_ "github.com/lib/pq"
	"source-pricing/models"
	"os"
	"github.com/gorilla/mux"
)

func PostgresClient() *sql.DB {

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func RedisClient() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: "",
		DB: 0,
	})
	return client

}




func GetProductDetailHandler(res http.ResponseWriter, req *http.Request) {
	
	
    //fetch parameters from request
	params := mux.Vars(req)
	product_id:=params["product_id"]
	store_id:=params["store_id"]
	redisKey := product_id + "_" + store_id
	
	//creating redis connection
	redisClient := RedisClient()
	pong, err := redisClient.Ping().Result()
	log.Println("Redis ping", pong, err)

    
	redisSearchResult, err := redisClient.Get(redisKey).Result()
	res.Header().Set("Content-Type", "application/json")
	
	//return the value from the redis if key found otherwise return from the db and update the cache
	if redisSearchResult != "" {
		fmt.Println("from redis result")
		var product models.Product
		json.Unmarshal([]byte(redisSearchResult), &product)
	
		jsonProductData, err_marshal := json.Marshal(&product)
		log.Println(jsonProductData,err_marshal)
        
	    json.NewEncoder(res).Encode(product)
	} else {
	db := PostgresClient()
    var product models.Product
	sqlStatement := `SELECT price, mrp FROM pricing_domain_product where retail_outlet_id=$1 and cms_product_id=$2`

	PID, err_product := strconv.Atoi(product_id)
	log.Println(err_product)

	SID,err_store:=strconv.Atoi(store_id)
	log.Println(err_store)

	row := db.QueryRow(sqlStatement, SID, PID)
	err:= row.Scan(&product.Price, &product.MRP)
	log.Println(err)

	redisValue, err_marshal := json.Marshal(&product)
	log.Println(redisValue,err_marshal)

	err_redis :=redisClient.Set(redisKey,string(redisValue),0 ).Err()
	log.Println(err_redis)

	defer db.Close()
    json.NewEncoder(res).Encode(product)
	}
 }