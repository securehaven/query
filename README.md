# SecureHaven: Query


## Getting started

Install the module by running

```
go get github.com/securehaven/query
```


## Usage

```go
import ("github.com/securehaven/query")

func main() {
	parser := query.NewParser(
		map[string]query.ParseFunc{
			"id":         query.ParseInt(0, 0),
			"first_name": query.ParseString,
			"last_name":  query.ParseString,
		},
	)

	queryValues, _ := url.ParseQuery("limit=10&offset=0&sort=id:asc&select=first_name,last_name&id=gt:1")
	query, err := parser.Parse(queryValues)

	if err != nil {
		fmt.Printf("failed to parse some value: %v", err)
	}

	fmt.Printf("%+v", query)
	// {Limit:10 Offset:0 Select:[first_name last_name] Sortings:[{Field:id Order:asc}] Filterings:[{Field:id Filter:gt Value:1}]}
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
