# Searcher 🔍
This project aims to provide a standardized approach for searching in different database engines within a single project. It addresses the challenge of enabling seamless search functionality across various databases without the need to implement different search methods for each database type. By using this project, developers can streamline the search process and ensure consistency across multiple projects.

## Status of the project ✅:
- [x] Query Translation to MongoDB
- [x] Query Translation to Elasticsearch
- [x] Support for 2 levels of query (e.g.((x=1) AND (y=2 OR (z=3))))
- [x] Support for date range queries
- [x] Basic logic operators support "AND" and "OR"
- [x] Basic operators support ">", "<", "<=", ">=" and "!="
- [ ] Query Translation to MySQL
- [ ] Query Translation to PostgreSQL
- [ ] Support for multiple levels of query (n levels of depth)
- [ ] Support for number range queries
- [ ] Super filters work with all basic operators ">", "<", "<=", ">=" and "!="


## ¿How to use this project?
We will do it in step by step explanation below that consists in first initialize the searcher, define the permitted query fields for every entity and finishing using the QueryTranslator for get a query in our desired Database Engine.

### 1. Initialize the QueryTranslator
For this step we recommend you to use dependencies inversion inside of your repository for the QueryTranslator.

- First this is how my repository is builded using the QueryTranslator
    ```go
    package repository

    import "errors"

    type MyRepository struct {
        QueryTranslator *QueryTranslator
    }

    func NewMyRepository(queryTranslator *QueryTranslator) (*MyRepository, error) {
        if queryTranslator == nil {
            return nil, errors.New("nil query translator provided")
        }
        return &MyRepository{QueryTranslator: queryTranslator}, nil
    }
    ```
- Now this is how in my main repository I build the query translator and pass it to the repository
    ```go
    package main

    import (
        "github.com/solrac97gr/searcher"
        "repository"
    )
    func main() {
        queryTranslator, err := searcher.NewQueryTranslator()
        if err != nil {
            panic(err)
        }

        myRepo, err := repository.NewMyRepository(queryTranslator)
        if err != nil {
            panic(err)
        }
    }
    ```
- Let's use the query translator in our repository for search.
### 2. Define the permitted fields for a entity
For this we will use a file that looks like this:
```go
package repository

import "github.com/solrac97gr/searcher/domain/models"

const (
	Name        models.Field = "name"
	Address     models.Field = "address"
	District    models.Field = "district"
	PhoneNumber models.Field = "phone_number"
	Email       models.Field = "email"
	CreatedAt   models.Field = "created_at"
)

const (
	ValidClientsFieldEntityName = "clients"
)

var ValidClientsField = map[string]models.FieldMetaData{
	Name.String(): {
		Field: Name,
		Type:  models.String,
	},
	Address.String(): {
		Field: Address,
		Type:  models.String,
	},
	District.String(): {
		Field: District,
		Type:  models.String,
	},
	PhoneNumber.String(): {
		Field: PhoneNumber,
		Type:  models.String,
	},
	Email.String(): {
		Field: Email,
		Type:  models.String,
	},
	CreatedAt.String(): {
		Field: CreatedAt,
		Type:  models.Date,
	},
}

var ValidClientFieldSet = models.ValidFields{
	EntityName: ValidClientsFieldEntityName,
	Fields:     ValidClientsField,
}
```

And for use this file in our query translator we will use the following method.

```go
    package main

    import (
        "github.com/solrac97gr/searcher"
        "repository"
    )

    func main() {
        queryTranslator, err := searcher.NewQueryTranslator()
        if err != nil {
            panic(err)
        }

        // Here we add a validation set for the client entity
        queryTranslator.AddValidFieldsSet(ValidClientFieldSet)

        myRepo, err := repository.NewMyRepository(queryTranslator)
        if err != nil {
            panic(err)
        }
    }
```

### 3. Now we will use our QueryTranslator for generate a Query for Mongo Database engine:
```go
package repository

import (
	"context"
	"errors"

	"github.com/solrac97gr/searcher/domain/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	// Search must get the criteria that is the query for be translated
	// in this case I return an array of interfaces but you can return your own entity
	Search(userID string, criteria *models.Criteria) (total int64, result []interface{}, err error)
}

type MyRepository struct {
	Collection      mongo.Collection
	QueryTranslator *QueryTranslator
}

func NewMyRepository(queryTranslator *QueryTranslator) (*MyRepository, error) {
	if queryTranslator == nil {
		return nil, errors.New("nil query translator provided")
	}
	return &MyRepository{QueryTranslator: queryTranslator}, nil
}

func (r *MyRepository) Search(userID string, criteria *models.Criteria) (total int64, result []interface{}, err error) {
	if err := criteria.Validate(); err != nil {
		return total, nil, err
	}

	// The parameters are the following
	// - The valid map entity name: in my case I use the constant that I defined you can also pass a string
	// - The criteria: is the actual query to be converted in the engine that you want.
	// - The super filters: this filters are applied in the top of query and only accept the "=" operator, it works like (user_id = example123 AND (YOUR_CRITERIA_CONVERTED))
	query, err := r.QueryTranslator.ToMongo(ValidClientsFieldEntityName, *criteria, []models.SuperFilter{
		{
			Field: "user_id",
			Value: userID,
		},
	})
	if err != nil {
		return total, nil, err
	}

	// For the case of mongo you need to pass the filters and sort separately so you can use the methods for get they
	filters, err := query.GetFilters()
	if err != nil {
		return total, nil, err
	}
	sorts, err := query.GetSorts()
	if err != nil {
		return total, nil, err
	}

	opts := options.Find()
	opts.SetLimit(int64(criteria.Pagination.Limit))
	opts.SetSkip(int64(criteria.Pagination.Offset))
	opts.SetSort(sorts)

	res, err := r.Collection.Find(context.Background(), filters, opts)
	if err != nil {
		return total, nil, err
	}

	err = res.All(context.Background(), &result)
	if err != nil {
		return total, nil, err
	}

	total, err = r.Collection.CountDocuments(context.Background(), filters)
	if err != nil {
		return total, nil, err
	}

	return total, result, nil
}

```

## Example of how a Request to a Search endpoint using the package will look
> Here we are assuming that you are using our criteria structure for your endpoint body.
```bash
/v1/clients/search/
{
    "pagination": {
        "limit": 0,
        "offset": 0
    },
    "query": {
        "logical": "and",
        "filters": [
            {
                "conditions": [
                    {
                        "field": "district",
                        "operator": "=",
                        "value": "Miraflores"
                    },
                    {
                        "field": "email",
                        "operator": "=",
                        "value": "carlos-test@email.com"
                    }
                ],
                "logical": "and"
            },
            {
                "conditions": [
                    {
                        "field": "created_at",
                        "operator": ">=",
                        "value": "2020-01-18T18:16:00.000Z"
                    }
                ],
                "logical": "and"
            }
        ]
    }
}
```