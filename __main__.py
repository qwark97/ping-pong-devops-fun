import logging

from aiohttp.web import Application
from aiohttp.web import run_app
from aiohttp import web

app = Application()

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)
handler = logging.StreamHandler()
f = logging.Formatter('%(asctime)s - %(message)s')
handler.setFormatter(f)
logger.addHandler(handler)

visits = {}


def prepare_response(name):
    try:
        visits[name] += 1
    except KeyError:
        visits[name] = 1

    counter_resp = ('%dst time' if visits[name] == 1
                    else ('%dnd time' if visits[name] == 2
                          else ('%drd time' if visits[name] == 3
                                else '%dth time'))) % visits[name]
    resp = f"Welcome {name}, {counter_resp}"
    return resp


async def post_handler(request):
    body = await request.text()
    name = body.split(" ")[-1]
    logger.info("From caller: %s", body)
    resp = prepare_response(name)
    logger.info("Response: %s\n", resp)
    return web.Response(text=resp)


app.add_routes(
    [web.route('POST', '/ping', post_handler)]
)

run_app(app, host='127.0.0.1', port=5555)
