document.addEventListener("DOMContentLoaded", () => {

    const dataElement = document.getElementById("dashboard-data");

    if (!dataElement) {
        return;
    }

    const dashboard = JSON.parse(dataElement.textContent);

    createVolumeChart(dashboard.volume);

    createLevelChart(dashboard.levels);

    createServiceChart(dashboard.services);

});

function createVolumeChart(data) {

    const canvas = document.getElementById("volumeChart");

    if (!canvas) return;

    new Chart(canvas, {
        type: "line",

        data: {

            labels: data.map(v => v.hour),

            datasets: [{
                label: "Logs",

                data: data.map(v => v.count),

                borderWidth: 2,

                tension: 0.3,

                fill: false

            }]
        },

        options: {

            responsive: true,

            maintainAspectRatio: true,

            aspectRatio: 2,

        }
    });

}

function createLevelChart(data) {

    const canvas = document.getElementById("levelChart");

    if (!canvas) return;

    new Chart(canvas, {

        type: "pie",

        data: {

            labels: data.map(v => v.level),

            datasets: [{

                data: data.map(v => v.count)

            }]

        },

        options: {

            responsive: true,

            maintainAspectRatio: true,

            aspectRatio: 1,

        }

    });

}

function createServiceChart(data) {

    const canvas = document.getElementById("serviceChart");

    if (!canvas) return;

    new Chart(canvas, {

        type: "bar",

        data: {

            labels: data.map(v => v.service),

            datasets: [{

                label: "Logs",

                data: data.map(v => v.count)

            }]

        },

        options: {

            responsive: true,

            maintainAspectRatio: true,

            aspectRatio: 2,

        }

    });

}
