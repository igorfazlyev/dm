import { useState, useEffect } from 'react';
import { Plus, FileText, Calendar, AlertCircle, Upload, CheckCircle, Download } from 'lucide-react';
import api from '../../lib/api';
import t from '../../lib/translations';

const Studies = () => {
  const [studies, setStudies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showUploadModal, setShowUploadModal] = useState(false);
  const [selectedFile, setSelectedFile] = useState(null);
  const [uploadProgress, setUploadProgress] = useState(null);
  const [uploading, setUploading] = useState(false);
  const [formData, setFormData] = useState({
    modality: 'CBCT',
    study_date: new Date().toISOString().split('T')[0],
  });
  const [error, setError] = useState('');

  useEffect(() => {
    fetchStudies();
  }, []);

  // Poll for processing studies status
  useEffect(() => {
    const processingStudies = studies.filter(s => s.status === 'processing');
    
    if (processingStudies.length === 0) return;

    const interval = setInterval(async () => {
      for (const study of processingStudies) {
        try {
          const response = await api.get(`/studies/${study.id}/status`);
          if (response.data.status !== study.status) {
            fetchStudies();
            break;
          }
        } catch (error) {
          console.error('Failed to check study status:', error);
        }
      }
    }, 10000);

    return () => clearInterval(interval);
  }, [studies]);

  const fetchStudies = async () => {
    try {
      const response = await api.get('/patient/studies');
      setStudies(response.data);
    } catch (error) {
      console.error('Failed to fetch studies:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleFileSelect = (e) => {
    const file = e.target.files[0];
    if (file) {
      if (!file.name.toLowerCase().endsWith('.dcm')) {
        setError('Пожалуйста, выберите файл .dcm');
        return;
      }
      setSelectedFile(file);
      setError('');
    }
  };

  const handleCreateAndUpload = async (e) => {
    e.preventDefault();
    
    if (!selectedFile) {
      setError('Пожалуйста, выберите файл для загрузки');
      return;
    }

    setUploading(true);
    setError('');
    setUploadProgress('Создание исследования...');

    try {
      // Step 1: Create study
      const createResponse = await api.post('/studies', formData);
      const studyId = createResponse.data.id;
      
      // Step 2: Initialize upload
      setUploadProgress('Подготовка к загрузке...');
      await api.post(`/studies/${studyId}/upload/init`);
      
      // Step 3: Upload file
      setUploadProgress('Загрузка файла в Diagnocat...');
      const uploadFormData = new FormData();
      uploadFormData.append('file', selectedFile);
      
      await api.post(`/studies/${studyId}/upload`, uploadFormData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });

      setUploadProgress('✅ Готово! Анализ начался.');
      
      // Wait and close
      setTimeout(() => {
        setShowUploadModal(false);
        setSelectedFile(null);
        setUploadProgress(null);
        setFormData({
          modality: 'CBCT',
          study_date: new Date().toISOString().split('T')[0],
        });
        fetchStudies();
      }, 2000);

    } catch (err) {
      setError(err.response?.data?.error || 'Не удалось создать исследование');
      setUploadProgress(null);
    } finally {
      setUploading(false);
    }
  };

  const handleDownloadPDF = async (study) => {
    try {
      const response = await api.get(`/studies/${study.id}/pdf`, {
        responseType: 'blob',
      });
      
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `report_${study.id.slice(0, 8)}.pdf`);
      document.body.appendChild(link);
      link.click();
      link.parentNode.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      alert('Не удалось скачать отчет');
    }
  };

  const getStatusColor = (status) => {
    const colors = {
      created: 'bg-gray-100 text-gray-800',
      uploading: 'bg-blue-100 text-blue-800',
      processing: 'bg-yellow-100 text-yellow-800',
      completed: 'bg-green-100 text-green-800',
      failed: 'bg-red-100 text-red-800',
    };
    return colors[status] || colors.created;
  };

  if (loading) {
    return <div className="text-center py-12">{t.common.loading}</div>;
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{t.patient.studies.title}</h1>
          <p className="text-gray-600 mt-2">{t.patient.studies.subtitle}</p>
        </div>
        <button
          onClick={() => setShowUploadModal(true)}
          className="btn-primary flex items-center gap-2"
        >
          <Plus className="w-5 h-5" />
          {t.patient.studies.newStudy}
        </button>
      </div>

      {studies.length === 0 ? (
        <div className="card text-center py-12">
          <FileText className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">{t.patient.studies.noStudies}</h3>
          <p className="text-gray-600 mb-4">{t.patient.studies.noStudiesText}</p>
          <button
            onClick={() => setShowUploadModal(true)}
            className="btn-primary inline-flex items-center gap-2"
          >
            <Plus className="w-5 h-5" />
            {t.patient.studies.createStudy}
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {studies.map((study) => (
            <div key={study.id} className="card hover:shadow-md transition-shadow">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-primary-100 rounded-lg">
                    <FileText className="w-5 h-5 text-primary-600" />
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">{study.modality || 'CBCT'}</p>
                    <p className="text-sm text-gray-500">
                      {study.study_date ? new Date(study.study_date).toLocaleDateString('ru-RU') : 'Нет даты'}
                    </p>
                  </div>
                </div>
                <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(study.status)}`}>
                  {t.studyStatus[study.status] || study.status}
                </span>
              </div>

              <div className="space-y-2 text-sm">
                <div className="flex items-center gap-2 text-gray-600">
                  <Calendar className="w-4 h-4" />
                  {t.patient.studies.createdDate}: {new Date(study.created_at).toLocaleDateString('ru-RU')}
                </div>
              </div>

              <div className="mt-4 pt-4 border-t border-gray-200 flex gap-2">
                {study.status === 'processing' && (
                  <div className="flex items-center gap-2 text-yellow-600 text-sm">
                    <div className="w-4 h-4 border-2 border-yellow-600 border-t-transparent rounded-full animate-spin" />
                    Обработка...
                  </div>
                )}
                {study.status === 'completed' && (
                  <button 
                    onClick={() => handleDownloadPDF(study)}
                    className="text-green-600 hover:text-green-700 font-medium text-sm flex items-center gap-2"
                  >
                    <Download className="w-4 h-4" />
                    Скачать отчет
                  </button>
                )}
                {study.status === 'failed' && (
                  <span className="text-red-600 text-sm">Ошибка обработки</span>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Create & Upload Modal - SINGLE STEP */}
      {showUploadModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-lg w-full p-6">
            <h2 className="text-xl font-bold text-gray-900 mb-4">Создать новое исследование</h2>

            {error && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-3 flex items-start gap-2 mb-4">
                <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
                <p className="text-sm text-red-600">{error}</p>
              </div>
            )}

            {uploadProgress && (
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 flex items-start gap-2 mb-4">
                {uploadProgress.includes('✅') ? (
                  <CheckCircle className="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" />
                ) : (
                  <div className="w-5 h-5 border-2 border-blue-600 border-t-transparent rounded-full animate-spin flex-shrink-0 mt-0.5" />
                )}
                <p className="text-sm text-blue-800">{uploadProgress}</p>
              </div>
            )}

            <form onSubmit={handleCreateAndUpload} className="space-y-4">
              {/* Modality */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Модальность
                </label>
                <select
                  value={formData.modality}
                  onChange={(e) => setFormData({ ...formData, modality: e.target.value })}
                  className="input"
                  disabled={uploading}
                >
                  <option value="CBCT">CBCT</option>
                  <option value="CT">CT</option>
                  <option value="Panoramic">Panoramic</option>
                </select>
              </div>

              {/* Study Date */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Дата исследования
                </label>
                <input
                  type="date"
                  value={formData.study_date}
                  onChange={(e) => setFormData({ ...formData, study_date: e.target.value })}
                  className="input"
                  required
                  disabled={uploading}
                />
              </div>

              {/* File Upload */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Загрузите DICOM файл *
                </label>
                <div className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-primary-400 transition-colors">
                  <input
                    type="file"
                    accept=".dcm"
                    onChange={handleFileSelect}
                    className="hidden"
                    id="dicom-upload"
                    disabled={uploading}
                    required
                  />
                  <label
                    htmlFor="dicom-upload"
                    className={`cursor-pointer flex flex-col items-center ${uploading ? 'opacity-50 cursor-not-allowed' : ''}`}
                  >
                    <Upload className="w-12 h-12 text-gray-400 mb-2" />
                    {selectedFile ? (
                      <>
                        <p className="text-sm font-medium text-gray-900">{selectedFile.name}</p>
                        <p className="text-xs text-gray-500 mt-1">
                          {(selectedFile.size / 1024 / 1024).toFixed(2)} MB
                        </p>
                      </>
                    ) : (
                      <>
                        <p className="text-sm font-medium text-gray-700">Нажмите для выбора файла</p>
                        <p className="text-xs text-gray-500 mt-1">или перетащите файл сюда</p>
                        <p className="text-xs text-gray-400 mt-2">Поддерживаются только .dcm файлы</p>
                      </>
                    )}
                  </label>
                </div>
              </div>

              {/* Buttons */}
              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    if (!uploading) {
                      setShowUploadModal(false);
                      setSelectedFile(null);
                      setError('');
                      setUploadProgress(null);
                      setFormData({
                        modality: 'CBCT',
                        study_date: new Date().toISOString().split('T')[0],
                      });
                    }
                  }}
                  className="btn-secondary flex-1"
                  disabled={uploading}
                >
                  Отмена
                </button>
                <button
                  type="submit"
                  className="btn-primary flex-1 flex items-center justify-center gap-2"
                  disabled={uploading || !selectedFile}
                >
                  {uploading ? (
                    <>
                      <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                      Загрузка...
                    </>
                  ) : (
                    <>
                      <Upload className="w-4 h-4" />
                      Создать исследование
                    </>
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Studies;
