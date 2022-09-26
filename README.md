# ITSFood-Commerce

A make-to-order food marketplace app powered by Gin & PostgreSQL.


 ITS Food Commerce API v1 Documentation
--------------------------------------------------------------------------------

================================================================================

 CONSENSUS
	>>> URL Prefix: /v1/
	>>> Response Format: JSON
================================================================================

 CONTENT OUTLINES
 
 ^ Collection
 ^ Parameter
 ^ Query strings allowed
 ^ Default Response
 ^ Endpoints (& examples) with request & response payloads
________________________________________________________________________________

 AUTHENTICATION 

 All requests must have token included in the 'Authorization' HTTP header.

 The token is retrieved everytime user logs in the system.

 Login:

 POST /login/user  --> for users
 	Request payload: {
		username: (string)
		password: (string)
 	}

	If success:
	Response payload: {
		data: {
			'user_data': {
				id,
				name,
				email,
				phone,
				unit_id,
				unit_name,
				group_id,
				group_name
			},
			'api_token': <api_token>
		},
		status: 'success',
		message: 'Login successful',
		description: 'The username and password correct. Get the 

authentication token from the data attribute.'
	}

	If passsword is wrong:
	Response payload: {
		data: '',
		status; 'failed',
		message: 'Wrong password.',
		description: 'The username & password entered did not match.'
	}

	If username does not exist:
	Response payload: {
		data: '',
		status: 'failed',
		message: 'Username entered does not exist in the system.',
		description: 'Username entered does not exist in the system.'
	}

 POST /login/supplier  --> for suppliers
	Request payload: {
		username: (string)
		password: (string)
 	}

	If success:
	Response payload: {
		data: {
			'api_token': <api_token>,
			'supplier': <supplier_info>
		},
		status: 'success',
		message: 'Login successful',
		description: 'The username and password correct. Get the 

authentication token from the data attribute.'
	}

	If passsword is wrong:
	Response payload: {
		data: '',
		status; 'failed',
		message: 'Wrong password.',
		description: 'The username & password entered did not match.'
	}

	If username does not exist:
	Response payload: {
		data: '',
		status: 'failed',
		message: 'Username entered does not exist in the system.',
		description: 'Username entered does not exist in the system.'
	}
________________________________________________________________________________


 Collection:

 > Menus

 Parameter(s):

 # {id}

 Query strings allowed:
```sh		
 ~ vendorId			
 ~ filters ["category", "minOrderQty", "maxOrderQty", "preOrderDays", "preOrderHours"]
 ~ price ["min", "max"]
 ~ search <name, vendor_name>
 ~ length					
 ~ page
 ~ orderBy ["retail_price", "name", "random"]
 ~ sort (asc, desc)
 ```

 Response:
```sh
 - Status (success or failed)
 - Messages (Error messages)
 - Description (Error/Exception)
 - Result (The data being fetched)
 ```
 * When the data is a bulk of resources

 Endpoints:

 - All menus
```sh
	GET - /menus
```

	Result payload:
```sh
	{
		"data": [{
			id: (int),
			name: (string),
			description: (string),
			vendor_id: (int),
			vendor_name: (int),
			category: (string),
			retail_price: (int),
			wholesale_price: (int),
			pre_order_hours: (int),
			pre_order_days: (int),
			min_order_qty: (int),
			max_order_qty: (int),
			image: (string) ---> thumbnail
		}],
		totalRows: <rows_count>
	}
```

 - Details of a menu
```sh
	GET - /menus/:id/details
```

	Result payload:
```sh 
		"data": {
			id: (int),
			name: (string),
			description: (string),
			vendor_id: (int),
			vendor_name: (int),
			category: (string),
			retail_price: (int),
			wholesale_price: (int),
			pre_order_hours: (int),
			pre_order_days: (int),
			min_order_qty: (int),
			max_order_qty: (int),
			image: (string) ---> thumbnail
		}
```

 - Random 8 products with type of food
	GET - /products?type=food&order_by=random&number=8

	Response payload: {[
		product_id: (string),
		product_name: (string),
		product_price: (string),
		product_cogs: (string),
		pre_order_days: (string),
		pre_order_hours: (string),
		min_order_quantity: (string),
		supplier_id: (string)
		supplier_name: (string)
		product_image_url: (string) ---> thumbnail
	]}

 - Show the details of a product
	GET - /products/{product_id}/show

	Response payload: {[
		product_id: (string),
		product_name: (string),
		product_price: (string),
		product_cogs: (string),
		pre_order_days: (string),
		pre_order_hours: (string),
		min_order_quantity: (string),
		product_description: (string)
		supplier_id: (string)
		supplier_name: (string)
		product_image_url: (string) ---> detail
	]}

 - Get products of a supplier
	GET - /products?supplier_id=3

	Response payload: {[
		product_id: (string),
		product_name: (string),
		product_price: (string),
		product_cogs: (string),
		pre_order_days: (string),
		pre_order_hours: (string),
		min_order_quantity: (string),
		supplier_id: (string)
		supplier_name: (string)
		product_image_url: (string) ---> thumbnail
	]}

 - Get products of a tag
	GET - /products?tag_id=28

	Response payload: {[
		product_id: (string),
		product_name: (string),
		product_price: (string),
		product_cogs: (string),
		pre_order_days: (string),
		pre_order_hours: (string),
		min_order_quantity: (string),
		supplier_id: (string)
		supplier_name: (string)
		product_image_url: (string) ---> thumbnail
	]}

 - Get products that match the query search and pagination
	GET - /products?search=bakery&page=3&per_page=10

	Response payload: {[
		product_id: (string),
		product_name: (string),
		product_price: (string),
		product_cogs: (string),
		pre_order_days: (string),
		pre_order_hours: (string),
		min_order_quantity: (string),
		supplier_id: (string)
		supplier_name: (string)
		product_image_url: (string) ---> thumbnail
	]}

 - Deactivate a product
	GET /products/{product_id}/deactivate

	Response payload: {
		data: (empty string)
		status: (success or failed)
		message: (message from application)
	}


 - Activate a product
	GET /products/{product_id}/activate

	Response payload: {
		data: (empty string)
		status: (success or failed)
		message: (message from application)
	}
________________________________________________________________________________

 Collection:

 > Tags

 Parameter(s):

 #

 Query strings allowed:

 ~

 Default response:

 # Status (success or failed)
 # Message (Short message from application)
 # Description (Error/Exception)

 Endpoints:

 - Get all tags
	GET /tags
	
	Response payload: {[
		tag_id: (string)
		tag_name: (string)
	]}

 - Get tag  with specified ID
	GET /tags/{tag_id}
	
	Response payload: {[
		tag_id: (string)
		tag_name: (string)
	]}

________________________________________________________________________________

 Collection:

 > Cart

 Parameter(s):

 # {product_id}
 # {quantity}

 Query strings  allowed:

 ~

 Default response:

 # Status (success or failed)
 # Message (Short message from application)
 # Description (Error/Exception)

 Endpoints:

 - Add a product to cart
	POST /cart

	Request payload: {
		product_id: (string),
		product_name: (string),
		product_price: (string),
		product_qty: (string),
		product_cogs: (string),
		product_minimum_order_qty: (string),
		product_pre_order_days: (string),
		product_pre_order_hours: (string),
		product_supplier_id: (string),
		product_supplier_name: (string)
	}

 - Get total items in the cart
	GET /cart/total

	Response payload: {
		total_item: (string)
	}

 - Get contents of the cart
	GET /cart/content

	Response payload: {[
		row_id: (string),
		product_id: (string),
		product_name: (string),
		product_price: (string),
		product_qty: (string),
		product_cogs: (string),
		product_pre_order_days: (string),
		product_pre_order_hours: (string),
		product_min_order_qty: (string),
		product_supplier_id: (string),
		product_supplier_name: (string),
		subtotal: (string)
	]}

 - Get details of the cart
	GET /cart/detail

	Response payload: {
		cart_content: [
			row_id: (string),
			product_id: (string),
			product_name: (string),
			product_price: (string),
			product_qty: (string),
			product_cogs: (string),
			product_pre_order_days: (string),
			product_pre_order_hours: (string),
			product_supplier_id: (string),
			product_supplier_name: (string),
			subtotal (string)
		],
		num_of_products: (string),
		products_quantity: (string),
		total: (string)
	}

 - Update cart
	PUT /cart/update

	Request payload: {
		rowid: (string),
		product_qty: (string)
	}

 - Delete item in the cart
	PUT /cart/delete

	Request payload: {
		rowid: (string)
	}

 - Reset cart
	GET /cart/clear

	Response payload: default response

________________________________________________________________________________
 Collection

 > Suppliers

 Parameter(s):

 # {supplier_id}

 Query strings allowed:

 ~

 Default response:
 
 # Data (can be an empty string)
 # Status (success or failed)
 # Message (short message from application)
 
 Endpoints:

 - Get supplier's products data
	GET /suppliers/{supplier_id}/products

	Response payload: {
		id,
		name,
		price,
		cogs,
		min_order_qty,
		pre_order_days,
		pre_order_hours,
		status
	}

 - Save Firebase token to database
	POST suppliers/save-firebase-token
	
	Request payload: {
		firebase_token: (string)
	}

 - Get purchases data of the corresponding supplier grouped by order ID
	GET /suppliers/orders/{order_id} --> if the order_id is omited, then will fetch all order containing the purchases for the corresponding supplier, if not, then will fetch only purchases that exist in the order id provided

	Response payload: {
		order_id,
		user_name,
		unit_name,
		ordered_at,
		ordered_for,
		ordered_to,
		info,
		status
	}

 - Get purchases data of the corresponding supplier
	GET /suppliers/orders/{order_id}/purchases

	Response payload: {
		order_id,
		product_name,
		price,
		cogs,
		qty,
		extra_cost,
		status,
		ordered_at,
		ordered_for,
		ordered_to,
		info
	}
________________________________________________________________________________

 Collection

 > Orders

 Parameter(s):

 # {order_id}
 # {order_detail_id}
 # {supplier_id}

 Query strings allowed:

 ~

 Default response:

 # Data (can be an empty string)
 # Status (success or failed)
 # Message (short message from application)
 # Description (a longer message)

 Endpoints:
 
 - Add an order
	POST /orders

	Request payload: {
		alamat: (string)
		waktu: (string) --> must be in YYYY-MM-DD HH:MM:SS format
		tujuan: (string)
		aktivitas: (string) --> 'Rutin' or 'Pengembangan'
		dana: (string) --> 'Non-PNBP', 'BPPTNBH', 'APBNK', or 'Pribadi'
		rincian: (string)
	}

 - Get details of a purchase
	GET /orders/purchase/{order_detail_id}/details

	Response payload: {
		data: {
			order_id,
			product_name,
			product_price,
			product_qty,
			subtotal,
			product_image,
			purchase_status
		},
		status: (success or failed),
		message: (message  from application)
		description: (longer message)
	}

 - Proceed purchase by supplier
	GET orders/purchase/{order_details_id}/proceed/{supplier_id}

	Response payload: {
		data: (empty string)
		status: (success or failed)
		message: (message from application)
		description: (longer message)
	}