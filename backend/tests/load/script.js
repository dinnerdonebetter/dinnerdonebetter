import http from 'k6/http';
import { check, sleep } from 'k6';

const targetDomain = "dinnerdonebetter.dev"
const baseURL = `https://api.${targetDomain}`
const cookieName = "ddb_api_cookie"

export const options = {
  // A number specifying the number of VUs to run concurrently.
  vus: 1,
  // A string specifying the total duration of the test run.
  duration: '10s',
};

function buildHeaders() {
  return { 'Content-Type': 'application/json' }
}

function makeString(length) {
    let result = '';
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    const charactersLength = characters.length;
    let counter = 0;
    while (counter < length) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
      counter += 1;
    }
    return result;
}

class User {
  constructor() {
    this.username = makeString(24);
    this.password = makeString(32);
    this.emailAddress = `${makeString(32)}@${targetDomain}`;
  }
}

function register() {
  const user = new User();

  const data = {
    username: user.username,
    password: user.password,
    emailAddress: user.emailAddress,
  }

  const res = http.post(`${baseURL}/users`, JSON.stringify(data), {
    headers: buildHeaders(),
  });
  check(res, {
    'has status 201': (r) => r.status === 201,
  });

  return user
}

function login(user) {
  const data = {
    username: user.username,
    password: user.password,
  }

  const res = http.post(`${baseURL}/users/login`, JSON.stringify(data), {
    headers: buildHeaders(),
  });
  check(res, {
    'has status 202': (r) => r.status === 202,
    'has cookie': (r) => r.cookies[cookieName] !== null,
  });
}

// The function that defines VU logic.
//
// See https://grafana.com/docs/k6/latest/examples/get-started-with-k6/ to learn more
// about authoring k6 scripts.
//
export default function() {
  const user = register();
  login(user)
  sleep(1);
}
