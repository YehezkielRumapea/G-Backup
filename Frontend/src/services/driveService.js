import apiClient from './api';

const driveService = {
  /**
   * List files from Google Drive
   */
  async listFiles(remoteName, path = '/', searchTerm = '') {
    try {
      const response = await apiClient.get('/browser/list', {
        remote_name: remoteName,
        path: path,
        search_term: searchTerm
      });
      return response.data;
    } catch (error) {
      console.error('List files error:', error);
      throw error;
    }
  }
};

export default driveService;
