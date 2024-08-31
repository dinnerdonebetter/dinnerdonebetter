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

class User {
  constructor() {
    this.username = `k6_${makeString(24)}`
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
    'register returns 201': (r) => r.status === 201,
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
    'login returns 202': (r) => r.status === 202,
    'login returns cookie': (r) => r.cookies[cookieName] !== null,
  });

  const jar = http.cookieJar();
  jar.set(baseURL, cookieName, res.cookies[cookieName]);
}

function listRecipes() {
  const res = http.get(`${baseURL}/api/v1/recipes`, {
    headers: buildHeaders(),
  });
  check(res, {
    'listing recipes returns 200': (r) => r.status === 200,
  });
}

function createRecipe() {
  // needs at least two steps
  const recipe = {
    name: "",
    minimumEstimatedPortions: "",
    pluralPortionName: "",
    portionName: "",
    slug: "",
    steps: [
      // at least one instrument or vessel required
      {
        preparationID: "",
        products: [
          {
            name: "",
            type: "ingredient",
            minimumQuantity: 1.23,
          },
        ],
      },
    ],
    yieldsComponentType: "unspecified",
  }

  const res = http.post(`${baseURL}/api/v1/recipes`, JSON.stringify(recipe), {
    headers: buildHeaders(),
  });
  check(res, {
    'listing recipes returns 200': (r) => r.status === 200,
  });
}

// The function that defines VU logic.
//
// See https://grafana.com/docs/k6/latest/examples/get-started-with-k6/ to learn more
// about authoring k6 scripts.
//
export default function() {
  login(register());
  listRecipes();
  sleep(1);
}

// util functions

// https://stackoverflow.com/a/1349426
function makeString(length) {
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  const charactersLength = characters.length;
  let result = '';

  let counter = 0;
  while (counter < length) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
    counter += 1;
  }

  return result;
}
