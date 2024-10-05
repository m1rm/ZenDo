<nav class="navbar navbar-expand-md navbar-dark bg-dark">
    <div class="container-fluid">
        <a class="navbar-brand" href="/">ZenDo</a>
        <button
            class="navbar-toggler"
            type="button"
            data-bs-toggle="collapse"
            data-bs-target="#nav"
            aria-controls="nav"
            aria-expanded="false"
            aria-label="Toggle navigation"
        >
            <span class="navbar-toggler-icon"></span>
        </button>

        <div class="collapse navbar-collapse" id="nav">
            <ul class="navbar-nav me-auto mb-2 mb-md-0">
                <li class="nav-item">
                    <a class="nav-link" href="/">Home</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="/todos">Todos</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="/about">About</a>
                </li>
            </ul>
        </div>
    </div>
</nav>

<div class="container-fluid content">
    <slot />
</div>

<footer
    class="container-fluid d-flex flex-wrap justify-content-between align-items-center my-4 border-top footer"
>
    <div class="d-flex align-items-center">
        <span class="mb-3 mb-md-0 text-body-secondary">© 2024 ZenDo</span>
    </div>

    <ul class="nav col-md-4 justify-content-end list-unstyled d-flex">
        <li class="ms-3">
            <a class="text-body-secondary" href="/about">About</a>
        </li>
        <li class="ms-3">
            <a class="text-body-secondary" href="/imprint">Imprint</a>
        </li>
        <li class="ms-3">
            <a class="text-body-secondary" href="privacy">Privacy</a>
        </li>
    </ul>
</footer>

<style>
    .content {
        min-height: calc(100vh - 50px); /* Adjust based on footer height */
        padding-bottom: 50px; /* Adjust based on footer height */
    }

    .footer {
        position: fixed;
        bottom: 0;
        width: 100%;
        height: 50px; /* Adjust as needed */
        background-color: #f8f9fa;
        text-align: center;
        line-height: 50px; /* Adjust to vertically center the text */
    }
</style>
