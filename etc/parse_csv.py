import datetime
import time
uuids = {}
uuid_points = {}
time_stamps = {}

with open('C:\\Users\\Willi\\Downloads\\data-1535087475807.csv') as f:
    for line in f:
        l = line.split(',')
        ts = l[0][1:14]
        uuid = l[1][1:-1]
        if uuid not in uuids:
            uuids[uuid] = None
        price = l[2][1:-1]
        d = datetime.datetime.strptime(ts, '%Y-%m-%d %H')
        t = int(time.mktime(d.timetuple()))
        if t not in time_stamps:
            time_stamps[t] = {}
        time_stamps[t][uuid] = price
uuids = list(uuids.keys())

with open('parsed_data.csv', 'w+') as f:
    f.write('ts,' + ','.join(uuids) + '\n')
    last_price = {}
    for u in uuids:
        last_price[u] = 0
    for ts, price_uuids in time_stamps.items():
        prices = []
        for u in uuids:
            if u in price_uuids:
                prices.append(str(price_uuids[u]))
                last_price[u] = price_uuids[u]
            else:
                prices.append(str(last_price[u]))
        f.write(str(ts) + ',' + ','.join(prices) + '\n')


