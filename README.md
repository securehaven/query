# SecureHaven: Query


## Syntax


### Pagination

```
/users?offset=0&limit=10
```

### Select

```
/users?select=id,firstName,lastName
```

### Sort

```
/users?sort=createdAt:desc,firstName:asc
```

### Filter

- Less than: `lt`
- Greater than: `gt`
- Less than or equals: `lte`
- Greater than or equals: `gte`
- Like: `like`
- Not: `not`

```
/users?firstName=John&lastName=like:D%&age=gte:18,lte:35
```
