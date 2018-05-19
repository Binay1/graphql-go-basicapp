package main

import(
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"log"

	"github.com/graphql-go/graphql"
)

type person struct{
	Name string
	Email string
}

var data map[string]person

var persontype = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "person",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},

			"email": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var root = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "root",
		Fields: graphql.Fields{
			"person": &graphql.Field{
				Type: persontype,
				Args: graphql.FieldConfigArgument{
					"Name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error){
						name, isOk := p.Args["Name"].(string)
							if isOk{
								return data[name], nil	
							}
						return nil, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: root,
	},
)


func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	log.Print("Executing query\n")
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	return result
}

func Getemails(w http.ResponseWriter, r *http.Request) {
	result := executeQuery(r.URL.Query().Get("query"), schema)
	json.NewEncoder(w).Encode(result)
}

func main() {
	_ = importJSONDataFromFile("data.json", &data)
	http.HandleFunc("/email", Getemails)
	fmt.Println("Listening at Localhost:8080")
	http.ListenAndServe(":8080", nil)
}

//Helper function to get json data

func importJSONDataFromFile(fileName string, result interface{}) (isOK bool) {
	isOK = true
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print("Error:", err)
		isOK = false
	}
	err = json.Unmarshal(content, result)
	if err != nil {
		isOK = false
		fmt.Print("Error:", err)
	}
	return
}