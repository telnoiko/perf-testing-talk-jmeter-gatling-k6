import {
    randomString,
    uuidv4,
} from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const generateUserData = () => ({
    name: `name-${randomString(3)}`,
    email: `${randomString(3)}@${randomString(3)}.uk`,
    password: randomString(10),
});

export const generateTaskData = () => ({
    description: randomString(5),
    completed: Math.random() > 0.5,
});