import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { Trend } from 'k6/metrics';

// Custom metrics to track specific endpoint performance
const discoveryLatency = new Trend('latency_api_discovery');
const productLatency = new Trend('latency_products');
const healthLatency = new Trend('latency_health');

export const options = {
  // Load profile: Ramp up to 20 users, hold, then ramp down
  stages: [
    { duration: '10s', target: 20 },
    { duration: '20s', target: 20 },
    { duration: '10s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be under 500ms
    http_req_failed: ['rate<0.01'],   // Error rate should be less than 1%
  },
};

// Configuration from environment variables
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8000';
const API_KEY = __ENV.API_KEY || 'your-api-key-here';

export default function () {
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': API_KEY,
    },
  };

  // 1. Health Check
  group('Health Check', function () {
    const res = http.get(`${BASE_URL}/health`);
    healthLatency.add(res.timings.duration);
    check(res, {
      'health: status is 200': (r) => r.status === 200,
      'health: has valid body': (r) => r.status === 200 && r.json().status === 'OK',
    });
  });

  // 2. API Discovery (New Dynamic Route)
  group('API Discovery', function () {
    const res = http.get(`${BASE_URL}/api`);
    discoveryLatency.add(res.timings.duration);
    check(res, {
      'discovery: status is 200': (r) => r.status === 200,
      'discovery: returns array': (r) => r.status === 200 && Array.isArray(r.json()),
    });
  });

  // 3. Products Listing
  group('Products API', function () {
    const res = http.get(`${BASE_URL}/api/produk`);
    productLatency.add(res.timings.duration);
    check(res, {
      'products: status is 200': (r) => r.status === 200,
      'products: returns array': (r) => r.status === 200 && Array.isArray(r.json()),
    });
  });

  // 4. Categories Listing
  group('Categories API', function () {
    const res = http.get(`${BASE_URL}/categories`);
    check(res, {
      'categories: status is 200': (r) => r.status === 200,
    });
  });

  // Realistic pacing between iterations
  sleep(Math.random() * 2 + 1);
}
