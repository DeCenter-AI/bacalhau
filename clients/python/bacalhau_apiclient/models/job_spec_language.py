# coding: utf-8

"""
    Bacalhau API

    This page is the reference of the Bacalhau REST API. Project docs are available at https://docs.bacalhau.org/. Find more information about Bacalhau at https://github.com/filecoin-project/bacalhau.  # noqa: E501

    OpenAPI spec version: 0.3.18.post4
    Contact: team@bacalhau.org
    Generated by: https://github.com/swagger-api/swagger-codegen.git
"""


import pprint
import re  # noqa: F401

import six

from bacalhau_apiclient.configuration import Configuration


class JobSpecLanguage(object):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    """
    Attributes:
      swagger_types (dict): The key is attribute name
                            and the value is attribute type.
      attribute_map (dict): The key is attribute name
                            and the value is json key in definition.
    """
    swagger_types = {
        'command': 'str',
        'deterministic_execution': 'bool',
        'job_context': 'JobSpecLanguageJobContext',
        'language': 'str',
        'language_version': 'str',
        'program_path': 'str',
        'requirements_path': 'str'
    }

    attribute_map = {
        'command': 'Command',
        'deterministic_execution': 'DeterministicExecution',
        'job_context': 'JobContext',
        'language': 'Language',
        'language_version': 'LanguageVersion',
        'program_path': 'ProgramPath',
        'requirements_path': 'RequirementsPath'
    }

    def __init__(self, command=None, deterministic_execution=None, job_context=None, language=None, language_version=None, program_path=None, requirements_path=None, _configuration=None):  # noqa: E501
        """JobSpecLanguage - a model defined in Swagger"""  # noqa: E501
        if _configuration is None:
            _configuration = Configuration()
        self._configuration = _configuration

        self._command = None
        self._deterministic_execution = None
        self._job_context = None
        self._language = None
        self._language_version = None
        self._program_path = None
        self._requirements_path = None
        self.discriminator = None

        if command is not None:
            self.command = command
        if deterministic_execution is not None:
            self.deterministic_execution = deterministic_execution
        if job_context is not None:
            self.job_context = job_context
        if language is not None:
            self.language = language
        if language_version is not None:
            self.language_version = language_version
        if program_path is not None:
            self.program_path = program_path
        if requirements_path is not None:
            self.requirements_path = requirements_path

    @property
    def command(self):
        """Gets the command of this JobSpecLanguage.  # noqa: E501

        optional program specified on commandline, like python -c \"print(1+1)\"  # noqa: E501

        :return: The command of this JobSpecLanguage.  # noqa: E501
        :rtype: str
        """
        return self._command

    @command.setter
    def command(self, command):
        """Sets the command of this JobSpecLanguage.

        optional program specified on commandline, like python -c \"print(1+1)\"  # noqa: E501

        :param command: The command of this JobSpecLanguage.  # noqa: E501
        :type: str
        """

        self._command = command

    @property
    def deterministic_execution(self):
        """Gets the deterministic_execution of this JobSpecLanguage.  # noqa: E501

        must this job be run in a deterministic context?  # noqa: E501

        :return: The deterministic_execution of this JobSpecLanguage.  # noqa: E501
        :rtype: bool
        """
        return self._deterministic_execution

    @deterministic_execution.setter
    def deterministic_execution(self, deterministic_execution):
        """Sets the deterministic_execution of this JobSpecLanguage.

        must this job be run in a deterministic context?  # noqa: E501

        :param deterministic_execution: The deterministic_execution of this JobSpecLanguage.  # noqa: E501
        :type: bool
        """

        self._deterministic_execution = deterministic_execution

    @property
    def job_context(self):
        """Gets the job_context of this JobSpecLanguage.  # noqa: E501


        :return: The job_context of this JobSpecLanguage.  # noqa: E501
        :rtype: JobSpecLanguageJobContext
        """
        return self._job_context

    @job_context.setter
    def job_context(self, job_context):
        """Sets the job_context of this JobSpecLanguage.


        :param job_context: The job_context of this JobSpecLanguage.  # noqa: E501
        :type: JobSpecLanguageJobContext
        """

        self._job_context = job_context

    @property
    def language(self):
        """Gets the language of this JobSpecLanguage.  # noqa: E501

        e.g. python  # noqa: E501

        :return: The language of this JobSpecLanguage.  # noqa: E501
        :rtype: str
        """
        return self._language

    @language.setter
    def language(self, language):
        """Sets the language of this JobSpecLanguage.

        e.g. python  # noqa: E501

        :param language: The language of this JobSpecLanguage.  # noqa: E501
        :type: str
        """

        self._language = language

    @property
    def language_version(self):
        """Gets the language_version of this JobSpecLanguage.  # noqa: E501

        e.g. 3.8  # noqa: E501

        :return: The language_version of this JobSpecLanguage.  # noqa: E501
        :rtype: str
        """
        return self._language_version

    @language_version.setter
    def language_version(self, language_version):
        """Sets the language_version of this JobSpecLanguage.

        e.g. 3.8  # noqa: E501

        :param language_version: The language_version of this JobSpecLanguage.  # noqa: E501
        :type: str
        """

        self._language_version = language_version

    @property
    def program_path(self):
        """Gets the program_path of this JobSpecLanguage.  # noqa: E501

        optional program path relative to the context dir. one of Command or ProgramPath must be specified  # noqa: E501

        :return: The program_path of this JobSpecLanguage.  # noqa: E501
        :rtype: str
        """
        return self._program_path

    @program_path.setter
    def program_path(self, program_path):
        """Sets the program_path of this JobSpecLanguage.

        optional program path relative to the context dir. one of Command or ProgramPath must be specified  # noqa: E501

        :param program_path: The program_path of this JobSpecLanguage.  # noqa: E501
        :type: str
        """

        self._program_path = program_path

    @property
    def requirements_path(self):
        """Gets the requirements_path of this JobSpecLanguage.  # noqa: E501

        optional requirements.txt (or equivalent) path relative to the context dir  # noqa: E501

        :return: The requirements_path of this JobSpecLanguage.  # noqa: E501
        :rtype: str
        """
        return self._requirements_path

    @requirements_path.setter
    def requirements_path(self, requirements_path):
        """Sets the requirements_path of this JobSpecLanguage.

        optional requirements.txt (or equivalent) path relative to the context dir  # noqa: E501

        :param requirements_path: The requirements_path of this JobSpecLanguage.  # noqa: E501
        :type: str
        """

        self._requirements_path = requirements_path

    def to_dict(self):
        """Returns the model properties as a dict"""
        result = {}

        for attr, _ in six.iteritems(self.swagger_types):
            value = getattr(self, attr)
            if isinstance(value, list):
                result[attr] = list(map(
                    lambda x: x.to_dict() if hasattr(x, "to_dict") else x,
                    value
                ))
            elif hasattr(value, "to_dict"):
                result[attr] = value.to_dict()
            elif isinstance(value, dict):
                result[attr] = dict(map(
                    lambda item: (item[0], item[1].to_dict())
                    if hasattr(item[1], "to_dict") else item,
                    value.items()
                ))
            else:
                result[attr] = value
        if issubclass(JobSpecLanguage, dict):
            for key, value in self.items():
                result[key] = value

        return result

    def to_str(self):
        """Returns the string representation of the model"""
        return pprint.pformat(self.to_dict())

    def __repr__(self):
        """For `print` and `pprint`"""
        return self.to_str()

    def __eq__(self, other):
        """Returns true if both objects are equal"""
        if not isinstance(other, JobSpecLanguage):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, JobSpecLanguage):
            return True

        return self.to_dict() != other.to_dict()
