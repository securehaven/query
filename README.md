# SecureHaven: Query


## Getting started

Install the module by running

```
go get github.com/securehaven/query
```


## Usage

```go
package main

import ("github.com/securehaven/query")

type exampleUser struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type examplePost struct {
	Id     int         `json:"id"`
	Title  string      `json:"title"`
	Author exampleUser `json:"author"`
}

var parser = query.MustParser(query.NewParser[examplePost]())

func main() {
	queryValues, _ := url.ParseQuery("limit=10&offset=0&sort=id:asc&select=title,author.first_name,author.last_name&id=gt:1")
	q, err := parser.Parse(queryValues)

	if err != nil {
		fmt.Printf("failed to parse some value: %v", err)
	}

	fmt.Printf("%+v", query)
	// {Limit:10 Offset:0 Select:[title author.first_name author.last_name] Sortings:[{Field:id Order:asc}] Filterings:[{Field:id Filter:gt Value:1}]}
}
```


## Syntax


### Pagination

```
offset=<int>&limit=<int>
```

**Example**

```
/users?offset=0&limit=10
```

### Select

```
select=<field>,...
```

**Example**

```
/users?select=id,firstName,lastName
```

### Sort

Allowed order values: `asc`, `desc`, `asc_nulls_first`, `desc_nulls_first`, `asc_nulls_last`, `desc_nulls_last`

```
sort=<field>:<order>,...
```

**Example**

```
/users?sort=createdAt:desc,firstName:asc
```

### Filter

Allowed filter values: `eq`, `lt`, `gt`, `lte`, `gte`, `like`, `neq`

```
<field>=<value>&<field>=<filter>:<value>
```

> `<field>=<value>` and `<field>=eq:<value>` behave the same

**Example**

```
/users?firstName=John&lastName=like:D%&age=gte:18,lte:35
```
