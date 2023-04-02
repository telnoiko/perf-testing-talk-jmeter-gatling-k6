import {sleep} from 'k6';
import {createUser, login, logoutUser} from './user.js';
import {createTask, deleteTask, updateTask} from './task.js';
import {Rate} from 'k6/metrics';

export const options = {
    vus: 1,
    duration: '1s',
    thresholds: {
        'failed_task_update': ['rate<0.1'],
        'failed_user_create': ['rate<0.1'],
        'http_req_duration': ['p(95)<400']
    }
};
const updateTaskFailRate = new Rate('failed_task_update');
const createUserFailRate = new Rate('failed_user_create');

export function setup() {
    return createUser()
}

export default function (data) {
    const auth = login(data.user);
    const id = createTask(auth);
    updateTask(id, auth);
    deleteTask(id, auth);
}

export function teardown(data) {
    logoutUser(data.authHeaders);
}
