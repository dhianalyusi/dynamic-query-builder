# Dynamic Query Builder

### Install

```
https://github.com/dhianalyusi/dynamic-query-builder
```

### Use

```
import queryBuilder "github.com/dhianalyusi/dynamic-query-builder"

var dqb queryBuilder.DQB
query := dqb.Table("customers").Where(dqb.Or(
    dqb.NewExpression("name", ">", "value"),
    dqb.NewExpression("gender", ">", "m"),
)).Where("age > 20").Limit(10).Offset(0).Select(
    "id, name",
    ).Join("left join", "orders", "customer_id", "id")
```

```
r := dqb.Table("customers").Limit(10).Offset(0).Select(
    "name, babibu, orders.date",
).Join("left join", "orders", "customer_id", "id")
```

```
s := dqb.Table("customers").Select(
    "name, babibu, orders.date",
).Join("left join", "orders", "customer_id", "id").Order("id ASC")
```