import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 100,          // Simulate 100 users
    duration: '30s',  // Last for 30 seconds
};

let index = 39458; // Starting index value

export default function () {
    const url = 'http://localhost:8080/shorten'; // Replace with your actual Gateway host/port

    const payload = JSON.stringify({
        original_url: `https://www.ptt.cc/bbs/Gossiping/index${index}.html`, // Use the current index
    });

    const headers = {
        'Content-Type': 'application/json',
        // 'Authorization': 'Bearer YOUR_TOKEN', // Uncomment if authentication is needed
    };

    const res = http.post(url, payload, { headers });

    check(res, {
        'status is 200': (r) => r.status === 200 || r.status === 201,
    });

    index--; // Decrease the index by 1 after each request
    sleep(0.1); // Each VU waits 0.1 seconds before sending the next request
}
