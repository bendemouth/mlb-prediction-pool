import React from "react";
import { Route, Routes } from "react-router-dom";
import UserProfile from "../pages/UserProfile";
import Leaderboard from "../pages/Leaderboard";
import Predictions from "../pages/Predictions";
import Home from "../pages/Home";

function Navbar() {
    return (
        <React.Fragment>
        <nav className="navbar">
            <h1>MLB Prediction Pool</h1>
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="/predictions">Predictions</a></li>
                <li><a href="/leaderboard">Leaderboard</a></li>
                <li><a href="/profile">Profile</a></li>
            </ul>
        </nav>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/predictions" element={<Predictions />} />
          <Route path="/leaderboard" element={<Leaderboard />} />
          <Route path="/profile" element={<UserProfile />} />
          <Route path="/profile/:userId" element={<UserProfile />} />
        </Routes>   
        </React.Fragment>
        
    );
}

export default Navbar;