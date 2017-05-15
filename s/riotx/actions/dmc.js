import ObjectAssign from 'object-assign';
import swagger from '../../swagger';
import constants from '../../core/constants';

// APIは必須でサポートしなければならない URI
const DMC_URI = '/dmc';

export default {
  get: context => {
    // TODO: apiがpromiseであればもうちょいキレイに書き直せる。
    return new Promise((resolve, reject) => {
      const operationObject = swagger.client.spec.paths[DMC_URI].get;

      if (!operationObject || !swagger.client.apis.dmc || !swagger.client.apis.dmc[operationObject.operationId]) {
        return reject(new Error(`[fetch] ${swagger.client.url}; system entry point not found. (uri: ${DMC_URI})`));
      }

      const api = swagger.getApiByOperationID(operationObject.operationId);

      const token = context.getter(constants.GETTER_ENDPOINTS_ONE, context.getter(constants.GETTER_CURRENT)).token;

      api({/** TODO get only support. */}, {
        // TODO https://github.com/swagger-api/swagger-js/issues/1036 でやりたい
        // TODO component.jsと共通化したい
        requestInterceptor: (req) => {
          req.headers['Authorization'] = token;
          console.log('Interceptor(request):', req);
        },
        responseInterceptor: (res) => {
          console.log('Interceptor(response):', res);
        }
      }).then(res => {
        if (!res.ok) {
          throw new Error(`[fetch] ${res.url} error.`);
        }

        console.log(`[fetch] ${res.url} success.`);

        return {
          response: res.obj,
          operationObject
        };
      }).then(res => {
        context.commit(constants.MUTATION_DMC, res);
        const key = context.getter(constants.GETTER_CURRENT);
        const endpoint = context.getter(constants.GETTER_ENDPOINTS_ONE, key);
        context.commit(constants.MUTATION_ENDPOINTS_UPDATE, key, ObjectAssign({}, endpoint, res.response));
      }).then(() => {
        resolve();
      }).catch(err => {
        reject(err);
      });
    });
  },

  remove: context => {
    return Promise
      .resolve()
      .then(() => {
        context.commit(constants.MUTATION_DMC_REMOVE);
      });
  }
};
