export async function createComment(post_uuid, user_uuid) {
    const form = document.getElementById("create-comment-form");
    const formData = new FormData(form);

    const data = {};
    formData.forEach((value, key) => (data[key] = value));

    data.post_uuid = post_uuid;
    data.user_uuid = user_uuid;

    try {
        const response = await fetch("/api/post/createComment", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            alert("Commentaire ajouté avec succès!");
        } else {
            const error = await response.json();
            alert("Erreur lors de l'ajout du commentaire: " + error.message);
        }
    } catch (error) {
        console.error("Erreur lors de l'envoi du commentaire:", error);
        alert("Une erreur s'est produite lors de l'envoi du commentaire. Veuillez réessayer.");
    }
}