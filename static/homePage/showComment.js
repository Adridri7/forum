import { DisplayMessages } from "./displayMessage.js";

export async function fetchAllcomments() {

    const usersPostContainer = document.getElementById('users-post');
    const firstPostItem = usersPostContainer.querySelector('.message-item');
    const postUuid = firstPostItem.getAttribute('post_uuid');
    console.log("fetch comment: ", postUuid)

    try {
        const response = await fetch(`/api/post/fetchComment/${postUuid}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
        });


        if (!response.ok) {
            throw new Error('Erreur lors de la récupération des commentaires');
        }


        const comment = await response.json();
        console.log(comment)
        console.log(comment[0].content)

        comment.forEach(comment => {
            DisplayMessages(comment);
        });

    } catch (error) {
        console.error(error);
    }
}