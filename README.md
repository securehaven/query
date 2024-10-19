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

func main() {
	parser := query.NewParser(
		query.NewField("id", query.ParseInt(0, 0), func(d Data) any {
			return query.NewNull(d.Id, d.Id > 0)
		}),
		query.NewField("first_name", query.ParseString, func(d Data) any {
			return query.NewNull(d.FirstName, len(d.FirstName) > 0)
		}),
		query.NewField("last_name", query.ParseString, func(d Data) any {
			return query.NewNull(d.LastName, len(d.LastName) > 0)
		}),
)

	queryValues, _ := url.ParseQuery("limit=10&offset=0&sort=id:asc&select=first_name,last_name&id=gt:1")
	query, err := parser.Parse(queryValues)

	if err != nil {
		fmt.Printf("failed to parse some value: %v", err)
	}

	fmt.Printf("%+v", query)
	// {Limit:10 Offset:0 Select:[first_name last_name] Sortings:[{Field:id Order:asc}] Filterings:[{Field:id Filter:gt Value:1}]}

	filtered := q.Filter(Data{
		Id:        3,
		FirstName: "John",
		LastName:  "Doe",
	})

	fmt.Println(filtered)
	// Raw: map[first_name:{John true} last_name:{Doe true}]
	// JSON: {"first_name":"John","last_name":"Doe"}
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
