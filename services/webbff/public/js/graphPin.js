function updateGraphPinned(userId, projectId, value) {
  fetch(`/${userId}/projects/${projectId}/graphpinned/${value}/`, {
    method: 'GET',
  })
    .then(response => {
      if (!response.ok) {
        throw new Error('Failed to update pinning project');
      }
      window.location.href = window.location.href;
    })
    .catch(error => {
      console.error(error);
    });
}