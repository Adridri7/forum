import { UserInfo } from "./app.js";
import { fetchCommentDetails, fetchPostDetails } from "./dashboard.js";
import { DisplayMessages } from "./displayMessage.js";


const notifDot = document.getElementById('notification-dot');

export async function fetchNotifications(isNotif = false) {
    const response = await fetch('/api/post/notifications', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ user_uuid: UserInfo.user_uuid })
    });

    if (response.ok) {
        const notifications = await response.json();
        if (isNotif) {
            displayNotifications(notifications);
        } else if (notifications.length > 0) {
            notifDot.style.display = 'block';
        }
    } else {
        console.error('Error fetching notifications');
    }
}

async function displayNotifications(notifications) {
    let response
    console.log("dans display :", notifications);
    const notificationList = document.getElementById('users-post');
    notificationList.innerHTML = '';

    for (const notification of notifications) {
        if (notification.target_type === "post") {
            response = await fetchPostDetails(notification.reference_id);
            console.log(response);
            DisplayMessages(response[0], false, true);
        } else {
            response = await fetchCommentDetails(notification.reference_id);
            DisplayMessages(response[0], true, true);
        }

        const commentElement = document.createElement('div');
        commentElement.innerHTML = `
                    <div class="notif-item">
                        ${notification.username} : "${notification.action}" your  ${notification.target_type} :
                    </div>
                `;
        const targetMessageItem = document.querySelector(`[post_uuid="${response[0].comment_id || response[0].post_uuid}"]`);
        console.log("target", targetMessageItem);

        if (targetMessageItem) {
            const messageContent = targetMessageItem.querySelector('.message-content');
            console.log(messageContent)
            if (messageContent) {
                messageContent.prepend(commentElement);
            }
        }
    }
    MarkNotifsAsRead();
}

async function MarkNotifsAsRead() {
    console.log("l'id du user :", UserInfo.user_uuid)
    const response = await fetch('/api/post/notificationsRead', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ user_uuid: UserInfo.user_uuid })
    });

    if (response.ok) {
        notifDot.style.display = 'none';
    } else {
        console.error('Error fetching notifications');
    }
}