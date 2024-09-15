import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

function Profile() {
  const [profile, setProfile] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await axios.get('http://localhost:8080/profile', {
          headers: { Authorization: localStorage.getItem('token') }
        });
        setProfile(response.data);
      } catch (error) {
        console.error('Failed to fetch profile:', error);
        navigate('/login');
      }
    };

    fetchProfile();
  }, [navigate]);

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  if (!profile) return <div className="text-center mt-5">Loading...</div>;

  return (
    <div className="container mt-5">
      <div className="row justify-content-center">
        <div className="col-md-6">
          <div className="card">
            <div className="card-body">
              <h2 className="card-title text-center mb-4">Profile</h2>
              <p className="text-center">Welcome, {profile.user}!</p>
              <div className="d-grid">
                <button onClick={handleLogout} className="btn btn-danger">Logout</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Profile;
