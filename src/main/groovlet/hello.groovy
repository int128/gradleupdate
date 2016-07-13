import com.example.Item
import groovy.json.JsonBuilder

final item = new Item(id: 'foo', name: 'bar')
item.save()

final savedItem = Item.get('foo')
assert savedItem

response.contentType = 'application/json'
println new JsonBuilder(savedItem)
