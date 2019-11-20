

# Message Structure

# Level 0: Base Message
The base message is the lowest level message used to wrap messages
Base message has a action, msg, and an optional request_id field

```json
{
    "action": "<action>"
    "msg": {<action message>}
    "response_id":"<optional response id>"
}
```
The response id is used so the frontend can make a async request to the backend, the backend will responed with the
same response_id that was sent to it so the frontend knows what do with the request and will be of message action response

## Connect Message
The very first message that needs to be sent down the websocket is the connect message.
The connect message supplies the connect token recieved from the /login endpoint
{
    "action": "connect",
    "msg": {
        "token":"<token>"
     }
}

The srver will then response with a connect message that looks something like:

```json
{
    "action": "connect",
    "msg": {
        "success": true,
        "config": {
            "fav": {
                "stocks": [],
                "users": []
            },
            "settings": {
                "changePercent": true
            }
        },
        "token": "m28no3ajVRE8VNQh1vm55pNcJmYCdYh-",
        "uuid": "24"
    }
}
```
The uuid fild contains the connected user's uuid.


## Object Message
New Objects in the system are sent down the websocket connection via a Object Message.
Objects are essently anything that has state and are identified with a guid. They also are
labled with a type string.
When you fist login to the websocket the server will dump all of the objects for that user.
```json
{
    "action":"object"
    "msg":{
        "type":"<object type>",
        "uuid": "<uuid for object>",
        "object": { <object value> }
    }
}
```

### User Object
all the Users are represnted with a User object that contains their display name, if they are logged in, and the uuid of
their portfolio. Users are there to represent the physical person in the game.
```json
}
{
    "action": "object",
    "msg": {
        "type": "user",
        "uuid": "24",
        "object": {
            "display_name": "DisplayName",
            "active": true,
            "portfolio_uuid": "25"
        }
    }
}
```


### Portfolio Object
The Portfolio objct contains the assioated user uud, the uuid of the portfolio, the current net worth, the level of the
portfolio, and the wallet. Portfolios are used to repsent the players game piece.

```json
{
    "action": "object",
    "msg": {
        "type": "portfolio",
        "uuid": "289",
        "object": {
            "user_uuid": "288",
            "uuid": "289",
            "wallet": 1001062322,
            "net_worth": 1001073780,
            "level": 0
        }
    }
}
```

### Stock Object
```json
{
    "action": "object",
    "msg": {
        "type": "stock",
        "uuid": "19",
        "object": {
            "uuid": "19",
            "name": "Michal Scott Paper Company",
            "ticker_id": "SCOTT",
            "current_price": 50233,
            "open_shares": 100
        }
    }
}
```

### Ledger Object
Ledgers are used to tie the amount of stock owned to a portfolio.
Ledgers are created as the relationships between stocks and portfolio are created (trade).
So if the srver were to be reset recently, there should be no ledger. If the amount of stock that
a user owns goes to zero, the ledger remains. This is so we can use the same uuid to track the entire history
of a stock-portfolio relation.
```json
{
    "action": "object",
    "msg": {
        "type": "ledger",
        "uuid": "66",
        "object": {
            "uuid": "66",
            "portfolio_id": "29",
            "stock_id": "17",
            "amount": 0,
            "record_book": "67"
        }
    }
}
```

### Record Book Objects
Record books tell the history of a ledger. They also keep track of the current active records.
A buy record is a purchase of a stock that has not been sold. This is used when calculating
profit and tax of a sale. Users only get their record books.
```json
{
    "action": "object",
    "msg": {
        "type": "record_book",
        "uuid": "47",
        "object": {
            "uuid": "47",
            "ledger_uuid": "46",
            "portfolio_uuid": "25",
            "buy_records": [
                {
                    "RecordUuid": "48",
                    "AmountLeft": 5
                },
                {
                    "RecordUuid": "56",
                    "AmountLeft": 2
                }
            ]
        }
    }
}
```
When selling a stock, the active records will be removed in a FIFO fashion. So in the example above, if 6 stocks are
sold, the amount of profit will be based on what the first 5 stocks were bought at, and then what 1 of the second
stock was bought at. THe resulting buy records will look like:
```json
 "buy_records": [
    {
        "RecordUuid": "56",
        "AmountLeft": 1
    }
]
```

### Record Objects
A single record keeps track of a single trade and any finance information about that trade.
Users only get their record objects
```json
{
    "action": "object",
    "msg": {
        "type": "record_entry",
        "uuid": "32",
        "object": {
            "uuid": "32",
            "share_price": 7459,
            "share_count": 15,
            "time": "2018-10-23T22:42:08.727523Z",
            "book_uuid": "31",
            "fee": 0,
            "taxes": 0,
            "bonus": 0,
            "result": -111885
        }
    }
}
```

### Notifcaitons Object
Notifactions are used as a way to notify users of events in the system.
While most of this information is dupliated elsewhere, Notifications are desinged to be
the user faceing alert of this information. Think facebook notifcations.

notifications are nested much the same way messages and objects are.

```json
    "action": "object",
    "msg": {
        "type": "notification",
        "uuid": "49",
        "object": {
            "uuid": "49",
            "user_uuid": "24",
            "time": "2018-10-23T22:42:31.278964Z",
            "type": "<notifcaiton type>",
            "notification": {<notifcation object>},
            "seen": false
        }
    }
}
```

#### Trade Notifcation
Tells about trades
```json
    "action": "object",
    "msg": {
        "type": "notification",
        "uuid": "49",
        "object": {
            "uuid": "49",
            "user_uuid": "24",
            "time": "2018-10-23T22:42:31.278964Z",
            "type": "trade",
            "notification": {
                "amount": 25,
                "stock": "21",
                "success": true
            },
            "seen": false
        }
    }
}
```
#### New Item Notifcations
Tells about the recieve of an item
```json
{
    "action": "object",
    "msg": {
        "type": "notification",
        "uuid": "509",
        "object": {
            "uuid": "509",
            "user_uuid": "24",
            "time": "2018-10-28T12:54:35.541021033Z",
            "type": "new_item",
            "notification": {
                "item_type": "insider",
                "item_uuid": "508"
            },
            "seen": false
        }
    }
}
```

### Item Objects
A item in the system. Users only get their items sent to them
```json
    "action": "object",
    "msg": {
        "type": "item",
        "uuid": "508",
        "object": {
            "type": "insider",
            "portfolio_uuid": "25",
            "uuid": "508",
            "used": false
        }
    }
}
```

## Update Action
Updates are closely tied to objects in that they represent the update of the state of an object.
Each update has the uuid of the object, the type of object it is, and a list of fields that have changed
on that object. The list is given with the field and value of each object.  This allows for objects to be updated
without sending the whole new object

Stock Update Example:

```json
{
    "action": "update",
    "msg": {
        "type": "stock",
        "uuid": "10",
        "changes": [
            {
                "field": "current_price",
                "value": 113521
            }
        ]
    }
}
```
If a stock update is triggered also all the profolio's that own that stock will be deleivered.

## Delete Action
Certin objects, notifcaions and items, can also be deleted from the system.
Delets can be sent from either the client or the server.
The delete flow looks something like:
```json
Client Sends
{
    "action": "delete",
    "msg": {
        "type": "item",
        "uuid": "508"
    },
    "request_id": "delete-response"
}
Revieve Response from server
{
    "action": "response",
    "msg": {
        "success": true,
        "err": ""
    },
    "request_id": "delete-response"
}
Revieve Delete Call from server
{
    "action": "delete",
    "msg": {
        "uuid": "508",
        "type": "item"
    }
}
```
with the idea being that the client will not acutally delete anything until the sever sends the delete action.

## Trade and Prospect Action
Tradeing is how stocks are puchased. Sending the stock id and amount to be sold.
Postive amount denotes a buy, negitive denotes a sell.

```json
Send:
{
    "action": "trade",
    "msg": {
        "stock_id": "<stock id>",
        "amount": <amount>
    },
    "request_id": "trade-response"
}
Response:
{
    "action": "response",
    "msg": {
        "order": {
            "stock_id": "<stock id>",
            "portfolio": "<user portfolio id>",
            "amount": <amount>
        },
        "details": {
            "share_price": <share price>,
            "share_count": <amount>,
            "shares_valuere": <share_price * amount>,
            "tax": <taxes if sell>,
            "fees": <fees on trade>,
            "bonus": <bonus money on trade>,
            "result": <share_value - tax - fees + bonus>
        },
        "success": true
    },
    "request_id": "trade-response"
}

```

Trades can also be prospected to see what the result of a trade would be. This takes into account any effects/items/whatever
the portfolio currently has active.
If someone buys a stock that they dont all ready own, then evryone will recive a new ledger object, else a ledger update object
will go out. If the trade is success, all the connected clients that belong to that user will be updated.
g