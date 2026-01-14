import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import Layout from './components/Layout';

// Auth pages
import Login from './pages/Login';
import Register from './pages/Register';

// Patient pages
import PatientDashboard from './pages/patient/Dashboard';
import Studies from './pages/patient/Studies';
import Plans from './pages/patient/Plans';
import Offers from './pages/patient/Offers';
import Profile from './pages/patient/Profile';

// Clinic pages
import ClinicDashboard from './pages/clinic/Dashboard';
import ClinicProfile from './pages/clinic/ClinicProfile';
import Pricelist from './pages/clinic/Pricelist';
import ClinicOffers from './pages/clinic/ClinicOffers';
import Orders from './pages/clinic/Orders';

// Protected Route Component
const ProtectedRoute = ({ children, allowedRoles }) => {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  if (allowedRoles && !allowedRoles.includes(user.role)) {
    return <Navigate to="/" replace />;
  }

  return <Layout>{children}</Layout>;
};

// Public Route (redirect if already logged in)
const PublicRoute = ({ children }) => {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  if (user) {
    // Redirect based on role
    if (user.role === 'patient') {
      return <Navigate to="/patient" replace />;
    } else if (user.role === 'clinic_manager' || user.role === 'clinic_doctor') {
      return <Navigate to="/clinic" replace />;
    }
  }

  return children;
};

// Root redirect based on user role
const RootRedirect = () => {
  const { user } = useAuth();

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  if (user.role === 'patient') {
    return <Navigate to="/patient" replace />;
  } else if (user.role === 'clinic_manager' || user.role === 'clinic_doctor') {
    return <Navigate to="/clinic" replace />;
  }

  return <Navigate to="/login" replace />;
};

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          {/* Root */}
          <Route path="/" element={<RootRedirect />} />

          {/* Public routes */}
          <Route
            path="/login"
            element={
              <PublicRoute>
                <Login />
              </PublicRoute>
            }
          />
          <Route
            path="/register"
            element={
              <PublicRoute>
                <Register />
              </PublicRoute>
            }
          />

          {/* Patient routes */}
          <Route
            path="/patient"
            element={
              <ProtectedRoute allowedRoles={['patient']}>
                <PatientDashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/patient/studies"
            element={
              <ProtectedRoute allowedRoles={['patient']}>
                <Studies />
              </ProtectedRoute>
            }
          />
          <Route
            path="/patient/plans"
            element={
              <ProtectedRoute allowedRoles={['patient']}>
                <Plans />
              </ProtectedRoute>
            }
          />
          <Route
            path="/patient/offers"
            element={
              <ProtectedRoute allowedRoles={['patient']}>
                <Offers />
              </ProtectedRoute>
            }
          />
          <Route
            path="/patient/profile"
            element={
              <ProtectedRoute allowedRoles={['patient']}>
                <Profile />
              </ProtectedRoute>
            }
          />

          {/* Clinic routes */}
          <Route
            path="/clinic"
            element={
              <ProtectedRoute allowedRoles={['clinic_manager', 'clinic_doctor']}>
                <ClinicDashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/clinic/profile"
            element={
              <ProtectedRoute allowedRoles={['clinic_manager', 'clinic_doctor']}>
                <ClinicProfile />
              </ProtectedRoute>
            }
          />
          <Route
            path="/clinic/pricelist"
            element={
              <ProtectedRoute allowedRoles={['clinic_manager', 'clinic_doctor']}>
                <Pricelist />
              </ProtectedRoute>
            }
          />
          <Route
            path="/clinic/offers"
            element={
              <ProtectedRoute allowedRoles={['clinic_manager', 'clinic_doctor']}>
                <ClinicOffers />
              </ProtectedRoute>
            }
          />
          <Route
            path="/clinic/orders"
            element={
              <ProtectedRoute allowedRoles={['clinic_manager', 'clinic_doctor']}>
                <Orders />
              </ProtectedRoute>
            }
          />

          {/* 404 */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;
