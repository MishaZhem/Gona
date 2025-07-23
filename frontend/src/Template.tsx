import { Outlet } from "react-router-dom"

function Template() {

    return (
        <div>
            <h1>This is website for Gona</h1>
            <Outlet></Outlet>
        </div>
    )
}

export default Template
