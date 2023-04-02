import http from 'k6/http';
import {Rate} from 'k6/metrics';
import {generateUserData} from './seeder.js';

const baseUrl = __ENV.HOSTNAME
const urls = {
    createUser: `${baseUrl}/users`,
    loginUser: `${baseUrl}/users/login`,
    logoutUser: `${baseUrl}/users/logoutAll`,
};
const createUserFailRate = new Rate('failed form user create');
const params = {
    headers: {
        'Content-Type': 'application/json',
    },
};

function authHeaders(token) {
    const authParams = {
        headers: {
            authorization: `Bearer ${token}`
        }
    }
    Object.assign(authParams.headers, params.headers)
    return authParams;
}

export const createUser = () => {
    const generatedData = generateUserData();

    const res = http.post(urls.createUser, JSON.stringify(generatedData), params);
    createUserFailRate.add(res.status !== 201);

    const headers = authHeaders(res.json().token);

    return {user: generatedData, authHeaders: headers};

}

export const login = (user) => {
    const loginData = {
        email: user.email,
        password: user.password,
    }
    const res = http.post(urls.loginUser, JSON.stringify(loginData), params);
    return authHeaders(res.json().token);
}

export const logoutUser = (params) => {
     http.post(urls.logoutUser, null, params);
}
