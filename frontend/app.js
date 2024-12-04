document.addEventListener('DOMContentLoaded', function () {
  fetch('api/events')
    .then((response) => response.json())
    .then((data) => {
      if (data == null) {
        console.log('no data to fetch');
        const msgContainer = document.getElementById('message');
        const el = document.createElement('p');
        msgContainer.appendChild(el);
        el.innerText = 'there are no events yet';
        return;
      }
      console.log('Fetched events: ', data);
      const eventsContainer = document.getElementById('events');
      data.forEach((event) => {
        const eventElement = document.createElement('div');
        eventElement.innerText = `${event.CreatedAt} - ${event.Severity.String}: ${event.Message.String}`;
        eventsContainer.appendChild(eventElement);
      });
    })
    .catch((error) => console.error('Error fetching events:', error));
});
