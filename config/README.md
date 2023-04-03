Chux Parser Configuration

Chux Parser is a REST API responsible for parsing the configuration from YAML and handling JSON data. This README will guide you through configuring Chux Parser, following the points discussed in this thread.
Configuration File

Create a YAML configuration file (e.g., config.development.yaml) with the following structure:

yaml

logging:
  level: info
aws:
  bucketName: chux-crawler
  downloadPath: ~/projects/chux/chux-parser
auth:
  issuerUrl: https://dev-29752729.okta.com/oauth2/default
  tokenUrl: ""

dataStores:
  dataStore:
    mongo:
      target: mongo
      uri: mongodb://localhost:27017
      timeout: 30
      databaseName: exampleDb
      collectionName: exampleCollection
    redis:
      target: redis
      uri: redis://localhost:6379
      timeout: 10
      databaseName: ""
      collectionName: ""

This configuration file includes settings for logging, AWS, authentication, and data stores (MongoDB and Redis in this example).
Loading Configuration

Use the LoadConfig function to load the configuration from the YAML file:

go

config, err := LoadConfig("development")
if err != nil {
  log.Fatalf("Failed to load configuration: %v", err)
}

This function reads the configuration file, unmarshals it into a ParserConfig struct, and initializes the BizObjConfig struct. Error handling is built into the function to ensure proper loading and parsing of the configuration file.
Setting BizObjConfig

Set the BizObjConfig from ParserConfig by iterating over the DataStoreMap and filling the BizObjConfig struct:

go

for dsName, ds := range config.BizObjConfig.DataStores.DataStoreMap {
  // Do something with the data store name (dsName) and data store configuration (ds)
}

Passing BizObjConfig to WithConfig

Pass the BizObjConfig to the WithConfig function of your Product or Article structs:

go

product := NewProduct().WithConfig(config.BizObjConfig)

This way, the configuration settings for the data stores are preserved and used by the Product or Article structs.
Processing JSON Data

To process JSON data efficiently, use the readJSONObjects function to read JSON objects from a file and send them as strings to a channel. You can then handle the JSON strings and errors in the calling code, for example:

go

filePath := "path/to/your/json/file.json"
out := make(chan string)
errOut := make(chan error)

go readJSONObjects(filePath, out, errOut)

for {
  select {
  case jsonStr, ok := <-out:
    if !ok {
      out = nil // Set the channel to nil to stop checking it
    } else {
      // Process the JSON string (e.g., pass it to Product.SetState())
      fmt.Println("JSON Object:", jsonStr)
    }
  case err, ok := <-errOut:
    if !ok {
      errOut = nil // Set the channel to nil to stop checking it
    } else {
      // Handle the error (e.g., log it, exit the program, or take other appropriate action)
      fmt.Println("Error:", err)
    }
  }

  // Break the loop when both channels are closed and set to nil
  if out == nil && errOut == nil {
    break
  }
}

This approach allows you to process JSON objects one by one, as they are read from the file, which is more efficient