// console.log("Script Loaded")

async function loadMetrics() {
    // console.log("Loading Metrics")

    const response = await fetch("/metrics");
    const data = await response.json();

    // console.log(data);

    document.getElementById("completed").textContent = data.completed;
    document.getElementById("failed").textContent = data.failed;
    document.getElementById("queued").textContent = data.queued;
    document.getElementById("throughput").textContent = data.throughput.toFixed(2) + " jobs/s";
    document.getElementById("utilization").textContent = (data.utilization * 100).toFixed(0) + "%";
}


async function loadWorkers() {

    const response = await fetch("/workers");
    const workers = await response.json();
    const container = document.getElementById("workers");

    container.innerHTML = "";

    workers.forEach(worker => {
        container.innerHTML += `
        <div class="worker ${worker.state === "Busy" ? "busy" : "idle"}">
            <h3>Worker ${worker.id}</h3>
            <p>${worker.state === "Busy" ? "🟢 Busy" : "⚪ Idle"}</p>
            <p>${worker.jobId ? `Job ${worker.jobId}` : "No Job"}</p>
        </div>
        `;
    });
}

async function runBenchmark() {

    const workers = document.querySelector(
        'input[name="workers"]:checked'
    ).value;

    const jobs = document.querySelector(
        'input[name="jobs"]:checked'
    ).value;

    document.getElementById("benchmarkStatus").textContent =
        "Status: Running...";

    await fetch("/benchmark", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            workers: Number(workers),
            jobs: Number(jobs),
        }),
    });

    checkBenchmarkCompletion();
}

async function checkBenchmarkCompletion() {

    const interval = setInterval(async () => {

        const response = await fetch("/metrics");
        const metrics = await response.json();

        if (metrics.queued === 0) {

            clearInterval(interval);

            document.getElementById("benchmarkStatus").textContent =
                "Status: Completed";

            loadBenchmarkResults();
        }

    }, 500);

}

async function loadBenchmarkResults() {

    const response = await fetch("/benchmarks");
    const history = await response.json();

    const table = document.getElementById("benchmarkTable");

    table.innerHTML = "";

    history.forEach((result,index) => {

        table.innerHTML += `
            <tr>
                <td>${index+1}</td>
                <td>${result.workers}</td>
                <td>${result.jobs}</td>
                <td>${result.seconds.toFixed(2)}</td>
                <td>${result.throughput.toFixed(2)}</td>
                <td>${result.utilization.toFixed(1)}%</td>
                <td>${result.completed}</td>
                <td>${result.failed}</td>
            </tr>
        `;

    });

}

loadMetrics();
loadWorkers();
loadBenchmarkResults();
setInterval(()=> {
    loadMetrics();
    loadWorkers();
},1000);