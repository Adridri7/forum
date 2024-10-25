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
        } else if (notifications && notifications.length > 0) {
            notifDot.style.display = 'block';
        }
    } else {
        console.error('Error fetching notifications');
    }
}

async function displayNotifications(notifications) {
    let response;
    const notificationList = document.querySelector(`.users-post[data-section="notifications"]`);
    notificationList.innerHTML = '';
    notificationList.style.flexGrow = "0";

    if (!notifications || notifications.length == 0) {
        // Create a message element
        const noNotificationsMessage = document.createElement('div');
        noNotificationsMessage.textContent = "There are no notifications.";
        noNotificationsMessage.classList.add('no-notifications-message');
        notificationList.appendChild(noNotificationsMessage);
        return;
    }

    for (const notification of notifications) {
        if (notification.target_type === "post") {
            response = await fetchPostDetails(notification.reference_id);
            console.log(response);
            DisplayMessages(response[0], "notifications", false, true);
        } else {
            response = await fetchCommentDetails(notification.reference_id);
            DisplayMessages(response[0], "notifications", true, true);
        }

        const commentElement = document.createElement('div');

        const targetMessageItems = document.querySelectorAll(`[post_uuid="${response[0].comment_id || response[0].post_uuid}"]`);

        targetMessageItems.forEach(targetMessageItem => {
            const messageContent = targetMessageItem.querySelector('.message-content');
            if (messageContent) {
                const existingNotif = messageContent.querySelector(`[data-uuid="${notification.notification_id}"]`);
                if (!existingNotif) {
                    commentElement.innerHTML = `
                    <div class="notif-item" data-uuid="${notification.notification_id}">
                        ${notification.username} : "${notification.action}" your ${notification.target_type} :
                    </div>
                `;
                    messageContent.prepend(commentElement);
                }
            }
        });
    }
    MarkNotifsAsRead();
}

async function MarkNotifsAsRead() {
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