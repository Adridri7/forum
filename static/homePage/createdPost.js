import { fetchPosts } from "./app.js";

export async function createPost(event) {
    event.preventDefault();


    const messageInput = document.getElementById("message");
    const messageContent = messageInput.value;

    // Fonction pour extraire les hashtags
    function extractHashtags(content) {
        const words = content.split(' ');
        const hashtags = words.filter(word => word.startsWith('#')).map(word => word.substring(1));
        return hashtags;
    }

    const hashtags = extractHashtags(messageContent); // Extraire les hashtags

    const data = {
        content: messageContent,
        hashtags: hashtags // Ajouter les hashtags aux données envoyées
    };

    try {
        const response = await fetch("http://localhost:8080/api/post/createPost", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });
        console.log(data)
        if (response.ok) {
            alert("Post créé avec succès!");
            messageInput.value = '';
            fetchPosts();
            window.location.href = "/"
        } else {
            const error = await response.json();
            alert("Erreur lors de la création du post: " + error.message);
        }
    } catch (error) {
        console.error("Erreur lors de la création du post:", error);
    }
}