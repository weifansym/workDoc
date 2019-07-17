## json2csv
json2csv是一款用于将JSON数据转换成CSV格式文件的库。
```
const json2csv = require('json2csv');
const fs = require('fs');

const fields = ['car', 'price', 'color'];
const myCars = [
  {
    "car": "Audi",
    "price": 40000,
    "color": "blue"
  }, {
    "car": "BMW",
    "price": 35000,
    "color": "black"
  }, {
    "car": "Porsche",
    "price": 60000,
    "color": "green"
  }
];

let csv = json2csv({ data: myCars, fields: fields });
 
fs.writeFile('file.csv', csv, function(err) {
  if (err) throw err;
  console.log('file saved');
});
```
