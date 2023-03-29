import {sleep} from 'k6';
import {createUser, login, logoutUser} from './user.js';
import {createTask, deleteTask, updateTask} from './task.js';
import {Rate} from 'k6/metrics';

export const options = {
    vus: 1,
    duration: '1s',
    thresholds: {
        'failed form task update': ['rate<0.1'],
        'failed form user create': ['rate<0.1'],
        'http_req_duration': ['p(95)<400']
    }
};
const updateTaskFailRate = new Rate('failed form task update');
const createUserFailRate = new Rate('failed form user create');

export function setup() {
    return createUser()
}

// login, create, update, delete
export default function (data) {
    const auth = login(data.user);
    const id = createTask(auth);
    updateTask(id, auth);
    deleteTask(id, auth);
}

export function teardown(data) {
    console.log(`teardown: ${JSON.stringify(data)} `);
    logoutUser(data.authHeaders);
}
