import {
    randomString,
    uuidv4,
} from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const generateUserData = () => ({
    name: uuidv4(),
    email: `${randomString(8)}@example.com`,
    password: randomString(10),
});

export const generateTaskData = () => ({
    description: randomString(30),
    completed: Math.random() > 0.5,
});