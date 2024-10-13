# SecureHaven: Query


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

Allowed filter values: `eq`, `lt`, `gt`, `lte`, `gte`, `like`, `not`

```
<field>=<value>&<field>=<filter>:<value>
```

> `<field>=<value>` and `<field>=eq:<value>` behave the same

**Example**

```
/users?firstName=John&lastName=like:D%&age=gte:18,lte:35
```
