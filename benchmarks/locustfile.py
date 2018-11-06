import os
import time

from random import choice

import redis

from locust import Locust
from locust import events
from locust import HttpLocust
from locust import TaskSet
from locust import task

class LimiterBehavior(TaskSet):
    def on_start(self):
        pass

    @task(1)
    def get_tokens_http(self):
        response = self.client.get(
            "/API/v1/limiter/test/1000/1/1000/1/",
            name='/API/v1/limiter/test/1000/1/1000/1/'
        )

        assert response.status_code == 200, response.text


class WebsiteLimiter(HttpLocust):
    task_set = LimiterBehavior
    min_wait = 1000
    max_wait = 4000
