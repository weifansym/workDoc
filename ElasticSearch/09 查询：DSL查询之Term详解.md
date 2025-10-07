## Term查询引入

如前文所述，查询分基于文本查询和基于词项的查询:

<img width="912" height="449" alt="image" src="https://github.com/user-attachments/assets/d83f34c5-50fb-41d6-b20b-206bf0e3192c" />

本文主要讲基于词项的查询。

<img width="2362" height="1514" alt="image" src="https://github.com/user-attachments/assets/423e61c6-3ea9-4ac1-af6a-1f4089ffe840" />

## Term查询

> 很多比较常用，也不难，就是需要结合实例理解。这里综合官方文档的内容，我设计一个测试场景的数据，以覆盖所有例子。@pdai

准备数据

```bash
PUT /test-dsl-term-level
{
  "mappings": {
    "properties": {
      "name": {
        "type": "keyword"
      },
      "programming_languages": {
        "type": "keyword"
      },
      "required_matches": {
        "type": "long"
      }
    }
  }
}

POST /test-dsl-term-level/_bulk
{ "index": { "_id": 1 }}
{"name": "Jane Smith", "programming_languages": [ "c++", "java" ], "required_matches": 2}
{ "index": { "_id": 2 }}
{"name": "Jason Response", "programming_languages": [ "java", "php" ], "required_matches": 2}
{ "index": { "_id": 3 }}
{"name": "Dave Pdai", "programming_languages": [ "java", "c++", "php" ], "required_matches": 3, "remarks": "hello world"}
```

### 字段是否存在:exist

由于多种原因，文档字段的索引值可能不存在：

-   源JSON中的字段是null或\[\]
-   该字段已”index” : false在映射中设置
-   字段值的长度超出ignore\_above了映射中的设置
-   字段值格式错误，并且ignore\_malformed已在映射中定义

所以exist表示查找是否存在字段。

<img width="3534" height="1630" alt="image" src="https://github.com/user-attachments/assets/130be06a-d737-45cb-ab96-bdc831b69465" />

### id查询:ids

ids 即对id查找

```bash
GET /test-dsl-term-level/_search
{
  "query": {
    "ids": {
      "values": [3, 1]
    }
  }
}
```

<img width="3554" height="1626" alt="image" src="https://github.com/user-attachments/assets/f59cdd45-ea88-41d5-ad3e-03dd4b56e1d2" />

### 前缀:prefix

通过前缀查找某个字段

```bash
GET /test-dsl-term-level/_search
{
  "query": {
    "prefix": {
      "name": {
        "value": "Jan"
      }
    }
  }
}
```

<img width="3526" height="1660" alt="image" src="https://github.com/user-attachments/assets/07038c6c-994c-4785-bc03-65e6f94a4000" />

### 分词匹配:term

前文最常见的根据分词查询

```bash
GET /test-dsl-term-level/_search
{
  "query": {
    "term": {
      "programming_languages": "php"
    }
  }
}
```

<img width="3580" height="1646" alt="image" src="https://github.com/user-attachments/assets/c53c630e-75c6-421d-9a1d-ef2e271c876a" />

### 多个分词匹配:terms

按照读个分词term匹配，它们是or的关系

```bash
GET /test-dsl-term-level/_search
{
  "query": {
    "terms": {
      "programming_languages": ["php","c++"]
    }
  }
}
```

<img width="3550" height="1642" alt="image" src="https://github.com/user-attachments/assets/ac6a8c5c-7f20-4e99-a94b-43e0a20b12a7" />

### 按某个数字字段分词匹配:term set

设计这种方式查询的初衷是用文档中的数字字段动态匹配查询满足term的个数

```bash
GET /test-dsl-term-level/_search
{
  "query": {
    "terms_set": {
      "programming_languages": {
        "terms": [ "java", "php" ],
        "minimum_should_match_field": "required_matches"
      }
    }
  }
}
```

<img width="3450" height="1514" alt="image" src="https://github.com/user-attachments/assets/f02fcebe-95e1-4e84-92f7-5b6ed63513bc" />

### 通配符:wildcard

通配符匹配，比如`*`

```bash
GET /test-dsl-term-level/_search
{
  "query": {
    "wildcard": {
      "name": {
        "value": "D*ai",
        "boost": 1.0,
        "rewrite": "constant_score"
      }
    }
  }
}
```

<img width="3542" height="1652" alt="image" src="https://github.com/user-attachments/assets/8069c188-b1ff-4040-9666-d6680275e337" />

### 范围:range

常常被用在数字或者日期范围的查询

```bash
GET /test-dsl-term-level/_search
{
  "query": {
    "range": {
      "required_matches": {
        "gte": 3,
        "lte": 4
      }
    }
  }
}
```

<img width="3522" height="1642" alt="image" src="https://github.com/user-attachments/assets/7d27d31d-f4a7-4aae-9ceb-88aa617eb462" />

### 正则:regexp

通过\[正则表达式\]查询

以”Jan”开头的name字段

```bash
GET /test-dsl-term-level/_search
{
  "query": {
    "regexp": {
      "name": {
        "value": "Ja.*",
        "case_insensitive": true
      }
    }
  }
}
```

<img width="3540" height="1508" alt="image" src="https://github.com/user-attachments/assets/48d1d018-4247-425d-a214-33f675a4b7e3" />

### 模糊匹配:fuzzy

官方文档对模糊匹配：编辑距离是将一个术语转换为另一个术语所需的一个字符更改的次数。这些更改可以包括：

-   更改字符（box→ fox）
-   删除字符（black→ lack）
-   插入字符（sic→ sick）
-   转置两个相邻字符（act→ cat）

```bash
GET /test-dsl-term-level/_search
{
  "query": {
    "fuzzy": {
      "remarks": {
        "value": "hell"
      }
    }
  }
}
```

<img width="3430" height="1514" alt="image" src="https://github.com/user-attachments/assets/cfd353a4-fe00-465d-b36f-d9e5344cfe16" />

## 参考文章

[https://www.elastic.co/guide/en/elasticsearch/reference/current/term-level-queries.html](https://www.elastic.co/guide/en/elasticsearch/reference/current/term-level-queries.html)
