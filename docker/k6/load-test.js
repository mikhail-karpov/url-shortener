import http from 'k6/http';
import { check, sleep } from 'k6'

const BASE_URL = 'http://host.docker.internal:8080';

export const options = {
  stages: [
    { duration: '30s', target: 50 },
    { duration: '1m', target: 50 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.05'],
    http_req_duration: ['p(95)<100']
  },
  noConnectionReuse: true,
}

export default function () {

  let response = shortenUrl("https://google.com");
  check(response, {'shorten url success': r => r.status === 200});

  sleep(1);

  const id = response.json().id;
  response = getUrl(id);
  check(response, {'get url success': r => r.status === 200});
}

function shortenUrl(longUrl) {
  const url = `${BASE_URL}/api/v1/shorten`;
  const body = JSON.stringify({long_url: longUrl});
  const params = {
    headers: {'Content-Type': 'application/json'},
    type: 'shorten_url'
  }

  return  http.post(url, body, params);
}

function getUrl(id) {
  const params = {type: 'get_url'};
  return http.get(`${BASE_URL}/api/v1/${id}`, params);
}
