import http from 'k6/http';
import {Rate} from 'k6/metrics';
import {generateUserData} from './seeder.js';

const baseUrl = __ENV.HOSTNAME
const createUserFailRate = new Rate('failed_user_create');
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

    const res = http.post(`${baseUrl}/users`, JSON.stringify(generatedData), params);
    createUserFailRate.add(res.status !== 201);

    const headers = authHeaders(res.json().token);

    return {user: generatedData, authHeaders: headers};

}

export const login = (user) => {
    const loginData = {
        email: user.email,
        password: user.password,
    }
    const res = http.post(`${baseUrl}/users/login`, JSON.stringify(loginData), params);
    return authHeaders(res.json().token);
}

export const logoutUser = (params) => {
     http.post(`${baseUrl}/users/logoutAll`, null, params);
}
