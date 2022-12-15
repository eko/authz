import http from 'k6/http';
import { check } from 'k6';
import { htmlReport } from 'https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.2/index.js';

const baseURL = 'http://localhost:8080/v1';

export function setup() {
  // Load resources.
  console.log('1/4 loading resources...');
  for (let i = 0; i < 10000; i++) {
    const data = JSON.stringify({
      id: 'post.' + i,
      kind: 'post',
      value: '' + i,
    });

    http.post(baseURL + '/resources', data, {
      headers: {
        'Content-Type': 'application/json',
      },
    });
  }

  // Load policies.
  console.log('2/4 loading policies...');

  for (let i = 0; i < 1000; i++) {
    const min = 1;
    const max = 1000;
  
    const randomInteger = Math.floor(Math.random() * (max - min + 1)) + min;

    const associatedResources = ['post.' + randomInteger];

    const data = JSON.stringify({
      id: 'policy-post-' + i,
      resources: associatedResources,
      actions: ['create', 'read', 'update', 'delete'],
    });

    http.post(baseURL + '/policies', data, {
      headers: {
        'Content-Type': 'application/json',
      },
    });
  }

  // Load roles.
  console.log('3/4 loading roles...');

  for (let i = 0; i < 50; i++) {
    const associatedPolicies = [];
    for (let j = 0; j < i; j++) {
      associatedPolicies.push('policy-post-' + j);
    }

    const data = JSON.stringify({
      id: 'role-post-' + i,
      policies: associatedPolicies,
    });

    http.post(baseURL + '/roles', data, {
      headers: {
        'Content-Type': 'application/json',
      },
    });
  }

  // Load principals.
  console.log('4/4 loading principals...');

  for (let i = 0; i < 10000; i++) {
    const associatedRoles = [];
    for (let j = 0; j < 10; j++) {
      associatedRoles.push('role-post-' + j);
    }
    
    const data = JSON.stringify({
      id: 'test-principal-' + i,
      roles: associatedRoles,
    });

    http.post(baseURL + '/principals', data, {
      headers: {
        'Content-Type': 'application/json',
      },
    });
  }
}

export const options = {
  setupTimeout: '10m',
  duration: '30s',
  target: 100,
};

export default function () {
  const min = 1;
  const max = 10000;

  const randomInteger = Math.floor(Math.random() * (max - min + 1)) + min;

  const data = JSON.stringify({
    checks: [
      {
        principal: 'test-principal-' + randomInteger,
        resource_kind: 'post',
        resource_value: '' + randomInteger,
        action: randomInteger % 2 ? 'create' : 'unknown-action',
      },
    ],
  });

  const response = http.post(baseURL + '/check', data, {
    headers: {
      'Content-Type': 'application/json',
    },
  });

  check(response, { 'status was 200': (r) => r.status == 200 });
}

export function handleSummary(data) {
  return {
    'result.html': htmlReport(data),
    stdout: textSummary(data, { indent: ' ', enableColors: true }),
  };
}