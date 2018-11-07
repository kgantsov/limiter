import os
import time
import uuid

from random import choice

import redis

from locust import Locust
from locust import events
from locust import HttpLocust
from locust import TaskSet
from locust import task


class LimiterBehavior(TaskSet):
    def on_start(self):
        self.user_id = uuid.uuid1()

    @task(1)
    def get_tokens_http(self):
        response = self.client.get(
            '/API/v1/limiter/user_{}/1000/1/1000/1/'.format(self.user_id),
            name='/API/v1/limiter/test/1000/1/1000/1/'
        )

        assert response.status_code == 200, response.text


class WebsiteLimiter(HttpLocust):
    task_set = LimiterBehavior
    min_wait = 1000
    max_wait = 4000
